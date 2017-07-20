package main

import (
	"github.com/trasa/watchmud/server"
	"github.com/trasa/watchmud/web"
)

func main() {
	//server.Init()
	gameserver := server.NewGameServer()
	go gameserver.Start()

	web.Init(gameserver)
	web.Start()

	// tell everybody to quit
	//close(server.GlobalQuit)
	// TODO some sort of server.GameServerInstance.Shutdown() ?
}
