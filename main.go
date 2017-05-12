package main

import (
	"github.com/trasa/watchmud/server"
)



func main() {
	server.Init()
	go server.GameServerInstance.Run()
	server.ConnectHttpServer()
}
