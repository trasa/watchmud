package main

import (
	"flag"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/trasa/watchmud/db"
	"github.com/trasa/watchmud/rpc"
	"github.com/trasa/watchmud/server"
	"github.com/trasa/watchmud/serverconfig"
	"github.com/trasa/watchmud/web"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"os"
)

var (
	serverConfigFile = flag.String("serverConfig", "./worldfiles/server.yaml", "location of the server config file")
	doHelp           = flag.Bool("help", false, "Show Help")
	doHelpAlias      = flag.Bool("h", false, "Show Help")
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "\t%s [flags]\n", os.Args[0])
	fmt.Fprint(os.Stderr, "Flags:\n")
	flag.PrintDefaults()
}

func main() {
	flag.Parse()
	if *doHelp || *doHelpAlias {
		usage()
		os.Exit(2)
		return
	}

	config, serverConfigErr := readServerConfig(*serverConfigFile)
	if serverConfigErr != nil {
		fmt.Fprintln(os.Stderr, "Error reading server configuration file")
		fmt.Fprintln(os.Stderr, serverConfigErr)
		usage()
		os.Exit(3)
		return
	}

	// init logging
	logFile, err := initLogging(config)
	if logFile != nil {
		defer logFile.Close()
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error initializing logging")
		fmt.Fprintln(os.Stderr, err)
		os.Exit(4)
		return
	}

	if err := db.Init(config); err != nil {
		log.Fatal().
			Err(err).
			Msg("Failed to initialize database persistence")
	}

	gameserver, err := server.NewGameServer(config.WorldFilesDir)
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("Failed to start NewGameServer")
	}
	go gameserver.Run()

	// grpc server
	rpcServer := rpc.NewServer(gameserver)
	go rpcServer.Run(config.ServerPort)

	web.Start(config.WebPort)

	// tell everybody to quit
	//close(server.GlobalQuit)
	// TODO some sort of server.GameServerInstance.Shutdown() ?
}

func readServerConfig(configFileName string) (*serverconfig.Config, error) {
	fmt.Println("Reading Server Config file from", configFileName)
	// read the configuration file
	configFileData, err := ioutil.ReadFile(configFileName)
	if err != nil {
		return nil, err
	}
	// parse the configuration file
	serverConfig := serverconfig.Config{}
	if err := yaml.UnmarshalStrict(configFileData, &serverConfig); err != nil {
		return nil, err
	}

	return &serverConfig, nil
}

// initialize the zerolog setup, logging to both console and file specified
// returns the *os.File so that the caller can defer the close of the log
// file appropriately.
func initLogging(config *serverconfig.Config) (*os.File, error) {
	f, err := os.OpenFile(config.Log.File, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return f, err
	}
	wrt := io.MultiWriter(os.Stdout, f)
	if logLevel, err := zerolog.ParseLevel(config.Log.Level); err != nil {
		return f, err
	} else {
		zerolog.SetGlobalLevel(logLevel)
	}
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: wrt})
	log.Info().Msg("Logging initialized.")
	return f, nil
}
