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
		// TODO add in time
		case message := <-server.incomingMessageBuffer:
			server.handleMessage(message)
		}
	}
}

func (server *GameServer) handleMessage(message *Message) {
	log.Printf("server incoming message: %s", message.Body)
	switch messageType := message.Body["msg_type"]; messageType {
	case "login":
		log.Printf("login received: %s", message.Body)
		server.handleLogin(message)
	default:
		log.Printf("UNHANDLED messageType: %s, body %s", messageType, message.Body)
	}
}

func (server *GameServer) handleLogin(message *Message) {
	// todo authentication and stuff
	// is this connection already authenticated?
	// create a Player, put the Player in a Room, other World state stuff
	// return the result back to the Player
}
