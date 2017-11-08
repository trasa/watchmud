package server

import (
	"github.com/trasa/watchmud/client"
	"github.com/trasa/watchmud/gameserver"
	"github.com/trasa/watchmud/message"
	"github.com/trasa/watchmud/world"
	"log"
	"time"
)

// Game Server ticks every PULSE_INTERVAL time
//const PULSE_INTERVAL time.Duration = 100 * time.Millisecond // 0.1 seconds
const PULSE_INTERVAL time.Duration = 1 * time.Second // 1 second

// Mobs consider doing something once every PULSE_MOBILE time
const PULSE_MOBILE = 10 * time.Second

// zones reset based on their lifetime, check every pulse_zone (lifetime > pulse_zone)
const PULSE_ZONE = 1 * time.Minute

type GameServer struct {
	incomingBuffer chan *gameserver.HandlerParameter
	World          *world.World
	tickInterval   time.Duration
}

func NewGameServer() *GameServer {
	return &GameServer{
		incomingBuffer: make(chan *gameserver.HandlerParameter),
		World:          world.New(),
	}
}

//noinspection SpellCheckingInspection
func (gs *GameServer) Run() {
	// this is the loop that handles incoming requests
	// needs to be organized around PULSEs
	tstart := time.Now().UnixNano()
	ticker := time.NewTicker(PULSE_INTERVAL)
	pulse := PulseCount(0)

	for {
		select {
		case <-ticker.C:
			now := time.Now().UnixNano()
			deltaSeconds := float64(now-tstart) / 1000000000
			tstart = now
			pulse++
			gs.heartbeat(pulse, deltaSeconds)

			//log.Printf("duration is %f seconds", pulse.toDuration().Seconds())

			// do something every 5 seconds
			//if pulse.checkInterval(5 * time.Second) {
			//	log.Printf("5 seconds")
			//}
			//if pulse.checkInterval(1 * time.Minute) {
			//	log.Printf("1 minute")
			//}
		}
	}
}

// runs the heartbeat of the game. Use pulse to determine intervals
// between things (ex. reset zones every 15 minutes...)
// delta is the amount of time since the last heartbeat was run.
func (gs *GameServer) heartbeat(pulse PulseCount, delta float64) {
	//log.Printf("pulse %d hb %d", pulse, delta)
	// mobs, scripts, ...

	// pulse zone
	// (zone reset ...)
	if pulse.checkInterval(PULSE_ZONE) {
		log.Printf("pulse %d zone reset %s", pulse, time.Now())
		// TODO
		gs.World.DoZoneActivity()
	}

	// pulse mobs
	// (mobs walk around, initiate attack?)
	if pulse.checkInterval(PULSE_MOBILE) {
		log.Printf("pulse %d do mobs %s", pulse, time.Now())
		gs.World.DoMobileActivity()
	}

	// perform violence
	// do the attacking (players and mobs and everybody)

	// mud-hour ("player tick")
	// affect weather, regen ..

	// handle an incoming message if one exists
	// TODO tick time: figure out how many incoming messages we can handle
	// see issue #4
	// for now, just process until buffer is empty...

	// not really infinite as the method will return false if there was
	// nothing to do.
	//noinspection GoInfiniteFor
	for gs.processIncomingMessage() {
	}
}

// read a message off of incomingMessageBuffer and do it
// this doesn't block so if the buffer is empty, the method returns immediately
// If a message was procssed (even in error) return true.
// Otherwise return false.
func (gs *GameServer) processIncomingMessage() bool {
	received := false
	select {
	case msg := <-gs.incomingBuffer:
		received = true
		switch msg.Message.Inner.(type) {
		case *message.GameMessage_LoginRequest:
			gs.handleLogin(msg) // TODO error handling
		default:
			gs.World.HandleIncomingMessage(msg)
		}
		/*
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
					Response: message.NewUnsuccessfulResponse("error", "TRANSLATE_ERROR"),
					Error:    fmt.Sprintf("%s", msg.Request.(message.ErrorRequest).Error),
				})

			default:
				gs.World.HandleIncomingMessage(msg)
					}
		*/

	default:
		//log.Printf("incoming message buffer is empty")
	}
	return received
}

func (gs *GameServer) Receive(msg *gameserver.HandlerParameter) {
	gs.incomingBuffer <- msg
}

func (gs *GameServer) Logout(c client.Client, cause string) {
	gm, err := message.NewGameMessage(message.LogoutRequest{Cause: cause})
	if err != nil {
		log.Printf("Error creating GameMessage for LogoutRequest: %v", err)
	} else {
		gs.Receive(gameserver.NewHandlerParameter(c, gm))
	}
}

func (gs *GameServer) handleLogin(msg *gameserver.HandlerParameter) { // TODO error handling
	// is this connection already authenticated?
	// see if we can find an existing player ..
	if msg.Client.GetPlayer() != nil {
		// you've already got one
		msg.Client.Send(message.LoginResponse{
			Success:    false,
			ResultCode: "PLAYER_ALREADY_ATTACHED",
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
	player := NewClientPlayer(msg.Message.GetLoginRequest().PlayerName, msg.Client)
	msg.Client.SetPlayer(player)
	msg.Player = player

	// add player to world
	gs.World.AddPlayer(player)
	player.Send(message.LoginResponse{
		Success:    true,
		ResultCode: "OK",
		PlayerName: player.GetName(),
	})
}
