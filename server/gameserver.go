package server

import (
	"fmt"
	"github.com/trasa/watchmud/client"
	"github.com/trasa/watchmud/message"
	"github.com/trasa/watchmud/world"
	"log"
	"time"
)

const TICK_INTERVAL time.Duration = 100 * time.Millisecond // 0.1 seconds

type GameServer struct {
	incomingMessageBuffer chan *message.IncomingMessage
	World                 *world.World
	tickInterval          time.Duration
}

func NewGameServer() *GameServer {
	return &GameServer{
		incomingMessageBuffer: make(chan *message.IncomingMessage),
		World: world.New(),
	}
}

func (gs *GameServer) Run() {
	// this is the loop that handles incoming requests
	// needs to be organized around TICKs
	tstart := time.Now().UnixNano()
	ticker := time.NewTicker(TICK_INTERVAL)

	for {
		select {
		case <-ticker.C:
			now := time.Now().UnixNano()
			deltaSeconds := float64(now-tstart) / 1000000000
			tstart = now

			gs.heartbeat(deltaSeconds)
		}
	}
}

func (gs *GameServer) heartbeat(delta float64) {
	log.Printf("hb %d", delta)
	// mobs, scripts, ...

	// handle an incoming message if one exists
	// TODO tick time: figure out how many incoming messages we can handle
	gs.processIncomingMessageBuffer()
}

// read a message off of incomingMessageBuffer and do it
// this doesn't block so if the buffer is empty, the method returns immediately
func (gs *GameServer) processIncomingMessageBuffer() {
	select {
	case msg := <-gs.incomingMessageBuffer:
		switch messageType := msg.Request.GetMessageType(); messageType {
		case "login":
			// login is special case, handled by server first and then
			// sent down to world for further initialization
			gs.handleLogin(msg) // TODO error handling

			// TODO does the gameserver need to do anything on logout?

		case "error":
			// we received bad input, send response to client
			// rather than processing message
			msg.Client.Send(message.ErrorResponse{
				Response: message.Response{
					Successful:  false,
					MessageType: "error",
					ResultCode:  "TRANSLATE_ERROR",
				},
				Error: fmt.Sprintf("%s", msg.Request.(message.ErrorRequest).Error),
			})

		default:
			gs.World.HandleIncomingMessage(msg)
		}
	default:
		log.Printf("incoming message buffer is empty")
	}

}

func (gs *GameServer) Receive(message *message.IncomingMessage) {
	gs.incomingMessageBuffer <- message
}

func (gs *GameServer) Logout(c client.Client, cause string) {
	gs.Receive(message.New(c, message.LogoutRequest{
		Request: message.RequestBase{MessageType: "logout"},
		Cause:   cause,
	}))
}

func (gs *GameServer) handleLogin(msg *message.IncomingMessage) { // TODO error handling
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
	player := NewClientPlayer(msg.Request.(message.LoginRequest).PlayerName, msg.Client)
	msg.Client.SetPlayer(player)
	msg.Player = player

	// add player to world
	gs.World.AddPlayer(player)
	player.Send(message.LoginResponse{
		Response: message.Response{
			MessageType: "login_response",
			Successful:  true,
			ResultCode:  "OK",
		},
		Player: message.NewPlayerData(player.GetName()),
	})
}
