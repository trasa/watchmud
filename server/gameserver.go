package server

import (
	"github.com/trasa/watchmud/message"
	"github.com/trasa/watchmud/world"
)

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

func (server *GameServer) Start() {
	// this is the loop that handles incoming requests
	// needs to be organized around TICKs
	for {
		select {
		// TODO add in tick time
		case msg := <-server.incomingMessageBuffer:
			switch messageType := msg.Body["msg_type"]; messageType {
			case "login":
				// login is special case, handled by server first and then
				// sent down to world for further initialization
				server.handleLogin(msg) // TODO error handling

			default:
				server.World.HandleIncomingMessage(msg)
			}
		}
	}
}

func (server *GameServer) Receive(message *message.IncomingMessage) {
	server.incomingMessageBuffer <- message
}

func (server *GameServer) handleLogin(msg *message.IncomingMessage) { // TODO error handling
	// is this connection already authenticated?
	// see if we can find an existing player ..
	if msg.Client.GetPlayer() != nil {
		// you've already got one
		msg.Client.Send(message.LoginResponse{
			Response: message.Response{
				MessageType: "login_response",
				Successful:  false,
				ResultCode:  "PLAYER_ALREADY_ATTACHED",
			},
		})
		return
	}
	// what if player is logged in on a different client?
	// TODO
	/*
		p := FindPlayerByClient(message.Client)
		if p != nil {
			// already authenticated, can't login again
			// TODO
			// note that this isn't really working; the same username can log on twice
			// instead the old player should be kicked and the new player take over
			p.Send(LoginResponse{
				Response: Response{
					MessageType: "login_response",
					Successful:  false,
					ResultCode:  "ALREADY_AUTHENTICATED",
				},
			})
			return
		}
	*/

	// todo authentication and stuff
	player := NewClientPlayer(msg.Body["player_name"], msg.Client)
	msg.Client.SetPlayer(player)
	msg.Player = player

	// add player to world
	server.World.AddPlayer(player)
	player.Send(message.LoginResponse{
		Response: message.Response{
			MessageType: "login_response",
			Successful:  true,
			ResultCode:  "OK",
		},
		Player: message.NewPlayerData(player.GetName()),
	})
}
