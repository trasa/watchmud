package main

import (
	"flag"
	"github.com/trasa/watchmud/rpc"
	"github.com/trasa/watchmud/server"
	"github.com/trasa/watchmud/web"
	"log"
)

var (
	worldFilesDir = flag.String("worldFilesDir", "./worldfiles", "directory where the world files can be found")
)

func main() {

	gameserver, err := server.NewGameServer(*worldFilesDir)
	if err != nil {
		log.Fatalf("Failed to start NewGameServer: %v", err)
	}
	go gameserver.Run()

	// grpc server
	rpcServer := rpc.NewServer(gameserver)
	go rpcServer.Run()

	web.Start()

	// tell everybody to quit
	//close(server.GlobalQuit)
	// TODO some sort of server.GameServerInstance.Shutdown() ?
}
