package server

import "github.com/trasa/watchmud/message"

//var GameServerInstance *GameServer

type GameServer struct {
	incomingMessageBuffer chan *message.IncomingMessage
	World                 *World
}

func NewGameServer() *GameServer {
	return &GameServer{
		incomingMessageBuffer: make(chan *message.IncomingMessage),
		World: NewWorld(),
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
			server.World.handleIncomingMessage(message)
		}
	}
}

func (server *GameServer) Receive(message *message.IncomingMessage) {
	server.incomingMessageBuffer <- message
}
