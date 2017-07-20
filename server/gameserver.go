package server

import (
	"github.com/trasa/watchmud/message"
	"github.com/trasa/watchmud/world"
)

//var GameServerInstance *GameServer

type GameServer struct {
	incomingMessageBuffer chan *message.IncomingMessage
	World                 *world.World
}

func NewGameServer() *GameServer {
	return &GameServer{
		incomingMessageBuffer: make(chan *message.IncomingMessage),
		World: world.New(),
	}
}

// Initialize the game server
//func Init() {
//	GameServerInstance = newGameServer()
//}

func (server *GameServer) Start() {
	// this is the loop that handles incoming requests
	// needs to be organized around TICKs
	for {
		select {
		// TODO add in tick time
		case message := <-server.incomingMessageBuffer:
			server.World.HandleIncomingMessage(message)
		}
	}
}

func (server *GameServer) Receive(message *message.IncomingMessage) {
	server.incomingMessageBuffer <- message
}
