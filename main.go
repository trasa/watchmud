package main

import (
	"github.com/trasa/watchmud/server"
	"github.com/trasa/watchmud/web"
)

func main() {
	server.Init()
	go server.GameServerInstance.Run()
	web.Init(server.GameServerInstance)
	web.ConnectHttpServer()

	// tell everybody to quit
	//close(server.GlobalQuit)
	// TODO some sort of server.GameServerInstance.Shutdown() ?
}
