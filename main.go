package main

import (
	"flag"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/trasa/watchmud/db"
	"github.com/trasa/watchmud/rpc"
	"github.com/trasa/watchmud/server"
	"github.com/trasa/watchmud/web"
	"io"
	"github.com/rs/zerolog/log"
	"os"
)

var (
	worldFilesDir = flag.String("worldFilesDir", "./worldfiles", "directory where the world files can be found")
	serverPort    = flag.Int("serverPort", 10000, "Port to operate the gRPC server on")
	webPort       = flag.Int("webPort", 8888, "Port to operate the web server on")
	doHelp        = flag.Bool("help", false, "Show Help")
	doHelpAlias   = flag.Bool("h", false, "Show Help")
	logFile       = flag.String("logFile", "/var/log/watchmud/watchmud-server.log", "File to write server logs to")
	debug = flag.Bool("debug", false, "Set log level to debug")
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

	// init logging
	f, err := os.OpenFile(*logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic("error opening log file")
	}
	defer f.Close()
	wrt := io.MultiWriter(os.Stdout, f)
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: wrt})
	log.Info().Msg("Logging initialized.")

	if err := db.Init(); err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize database persistence")
	}

	gameserver, err := server.NewGameServer(*worldFilesDir)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to start NewGameServer")
	}
	go gameserver.Run()

	// grpc server
	rpcServer := rpc.NewServer(gameserver)
	go rpcServer.Run(*serverPort)

	web.Start(*webPort)

	// tell everybody to quit
	//close(server.GlobalQuit)
	// TODO some sort of server.GameServerInstance.Shutdown() ?
}
