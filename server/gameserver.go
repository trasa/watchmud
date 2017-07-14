package server

import (
	"log"
)

var GameServerInstance *GameServer

type GameServer struct {
	incomingMessageBuffer chan *IncomingMessage
	World                 *World
}

func newGameServer() *GameServer {
	return &GameServer{
		incomingMessageBuffer: make(chan *IncomingMessage),
		World: NewWorld(),
	}
}

// Initialize the game server
func Init() {
	GameServerInstance = newGameServer()
}

func (server *GameServer) Run() {

	StartAllClientDispatcher()

	// this is the loop that handles incoming requests
	// needs to be organized around TICKs
	for {
		select {
		// TODO add in tick time
		case message := <-server.incomingMessageBuffer:
			server.handleIncomingMessage(message)
		}
	}
}

func (server *GameServer) handleIncomingMessage(message *IncomingMessage) {
	// TODO this should be moved to world as well i think
	log.Printf("server incoming message: %s", message.Body)
	switch messageType := message.Body["msg_type"]; messageType {
	case "login":
		log.Printf("login received: %s", message.Body)
		server.World.HandleLogin(message)
	case "tell":
		log.Printf("tell: %s", message.Body)
		// TODO
		//server.handleTell(message)
	case "tell_all":
		log.Printf("Tell All: %s", message.Body)
		server.handleTellAll(message)
	default:
		log.Printf("UNHANDLED messageType: %s, body %s", messageType, message.Body)
	}
}

//func (server *GameServer) handleTell(message *IncomingMessage) {
//	from := message.Client.Player.Name
//	to := message.Body["to"]
//	msg := message.Body["message"]
//}

// Tell everybody in the game something.
// TODO this belongs somewhere else.
func (server *GameServer) handleTellAll(message *IncomingMessage) {
	if val, ok := message.Body["message"]; ok {
		// TODO need notification type
		SendToAllClients(TellAllResponse{
			Response: Response{
				MessageType: "tell_all_notification",
				Successful:  true,
				ResultCode:  val,
			},
			Sender: message.Client.Player.Name,
		})
	} else {
		message.Client.source <- TellAllResponse{
			Response: Response{
				MessageType: "tell_all_notification",
				Successful:  false,
				ResultCode:  "NO_MESSAGE",
			},
		}
	}
}
