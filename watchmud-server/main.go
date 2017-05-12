package main

var GameServerInstance *GameServer

func main() {
	GameServerInstance = newGameServer()
	go GameServerInstance.run()
	connectHttpServer()
}
