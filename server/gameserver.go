package server

import (
	"github.com/trasa/watchmud/world"
	"log"
)

var GameServerInstance *GameServer

type GameServer struct {
	incomingMessageBuffer chan *Message
	World                 *world.World
}

func newGameServer() *GameServer {
	return &GameServer{
		incomingMessageBuffer: make(chan *Message),
		World: world.NewWorld(),
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
	case "tell_all":
		log.Printf("Tell All: %s", message.Body)
		server.tellAll(message)
	default:
		log.Printf("UNHANDLED messageType: %s, body %s", messageType, message.Body)
	}
}

type Response struct {
	MessageType string `json:"msg_type"`
	Successful  bool   `json:"success"`
	ResultCode  string `json:"result_code"`
}

type LoginResponse struct {
	MessageType string        `json:"msg_type"`
	Successful  bool          `json:"success"`
	ResultCode  string        `json:"result_code"`
	Player      *world.Player `json:"player"`
}

// Authenticate, create a Player, put the Player in a Room,
// other World state stuff.
func (server *GameServer) handleLogin(message *Message) {
	// todo authentication and stuff
	// is this connection already authenticated?
	if message.Client.Player != nil {
		// already authenticated, can't login again
		message.Client.source <- LoginResponse{
			Successful: false,
			ResultCode: "ALREADY_AUTHENTICATED",
		}
		return
	}
	player := world.NewPlayer(message.Body["player_name"], message.Body["player_name"])
	message.Client.Player = player
	server.World.AddPlayer(player)
	message.Client.source <- LoginResponse{
		MessageType: "LoginResponse",
		Successful:  true,
		ResultCode:  "OK",
		Player:      player,
	}
}

// Tell everybody in the game something.
// TODO this belongs somewhere else.
func (server *GameServer) tellAll(message *Message) {
	if val, ok := message.Body["message"]; ok {
		// TODO need notification type
		Broadcaster <- Response{
			MessageType: "TellAllNotification",
			Successful:  true,
			ResultCode:  val,
		}
	} else {
		message.Client.source <- Response{
			MessageType: "TellAllResponse",
			Successful:  false,
			ResultCode:  "NO_MESSAGE",
		}
	}
}
