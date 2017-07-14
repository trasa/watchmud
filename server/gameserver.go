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
	log.Printf("server incoming message: %s", message.Body)
	switch messageType := message.Body["msg_type"]; messageType {
	case "login":
		log.Printf("login received: %s", message.Body)
		server.handleLogin(message)
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

// Authenticate, create a Player, put the Player in a Room,
// other World state stuff.
func (server *GameServer) handleLogin(message *IncomingMessage) {
	// todo authentication and stuff
	// is this connection already authenticated?
	if message.Client.Player != nil {
		// already authenticated, can't login again
		lr := LoginResponse{
			Response: Response{
				MessageType: "login_response",
				Successful:  false,
				ResultCode:  "ALREADY_AUTHENTICATED",
			},
		}
		message.Client.source <- lr
		log.Printf("login response %s", lr.MessageType)
		return
	}
	player := NewPlayer(message.Body["player_name"], message.Body["player_name"], message.Client)
	message.Client.Player = player
	server.World.AddPlayer(player)
	message.Client.source <- LoginResponse{
		Response: Response{
			MessageType: "login_response",
			Successful:  true,
			ResultCode:  "OK",
		},
		Player: player,
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
