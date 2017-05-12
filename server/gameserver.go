package server

import (
	"log"
)

var GameServerInstance *GameServer

type GameServer struct {
	incomingMessageBuffer chan *Message
}

func newGameServer() *GameServer {
	return &GameServer{
		incomingMessageBuffer: make(chan *Message),
	}
}

// Initialize the game server
func Init() {
	GameServerInstance = newGameServer()
}

func (server *GameServer) Run() {
	for {
		select {
		case message := <-server.incomingMessageBuffer:
			log.Printf("server processed incoming message: %s", message.Body)
		}
	}
}
