package main

import (
	"github.com/trasa/watchmud/server"
)

func main() {
	server.Init()
	go server.GameServerInstance.Run()
	server.ConnectHttpServer()

	// tell everybody to quit
	//close(server.GlobalQuit)
	// TODO some sort of server.GameServerInstance.Shutdown() ?
}
