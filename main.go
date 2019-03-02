package main

import (
	"flag"
	"fmt"
	"github.com/trasa/watchmud/db"
	"github.com/trasa/watchmud/rpc"
	"github.com/trasa/watchmud/server"
	"github.com/trasa/watchmud/web"
	"io"
	"log"
	"os"
)

var (
	worldFilesDir = flag.String("worldFilesDir", "./worldfiles", "directory where the world files can be found")
	serverPort    = flag.Int("serverPort", 10000, "Port to operate the gRPC server on")
	webPort       = flag.Int("webPort", 8888, "Port to operate the web server on")
	doHelp        = flag.Bool("help", false, "Show Help")
	doHelpAlias   = flag.Bool("h", false, "Show Help")
	logFile       = flag.String("logFile", "/var/log/watchmud/watchmud-server.log", "File to write server logs to")
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
		log.Fatalf("error opening log file: %v", err)
	}
	defer f.Close()
	wrt := io.MultiWriter(os.Stdout, f)
	log.SetOutput(wrt)
	log.Println("Logging initialized.")

	if err := db.Init(); err != nil {
		log.Fatalf("Failed to initialize database persistence: %v", err)
	}

	gameserver, err := server.NewGameServer(*worldFilesDir)
	if err != nil {
		log.Fatalf("Failed to start NewGameServer: %v", err)
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
