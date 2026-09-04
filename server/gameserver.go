package server

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/trasa/watchmud-message"
	"github.com/trasa/watchmud/client"
	"github.com/trasa/watchmud/gameserver"
	"github.com/trasa/watchmud/mudtime"
	"github.com/trasa/watchmud/world"
)

type GameServer struct {
	incomingBuffer chan *gameserver.HandlerParameter
	world          *world.World
	tickInterval   time.Duration
}

func New(w *world.World) *GameServer {
	return &GameServer{
		incomingBuffer: make(chan *gameserver.HandlerParameter),
		world:          w,
	}
}

func (gs *GameServer) Run(ctx context.Context) error {
	ticker := time.NewTicker(mudtime.PulseInterval)
	defer ticker.Stop()

	last := time.Now()
	var pulse mudtime.PulseCount

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case <-ticker.C:
			now := time.Now()
			delta := now.Sub(last)
			last = now
			pulse++
			gs.heartbeat(pulse, delta)
		}
	}
}

// runs the heartbeat of the game. Use pulse to determine intervals
// between things (ex. reset zones every 15 minutes...)
// delta is the amount of time since the last heartbeat was run.
func (gs *GameServer) heartbeat(pulse mudtime.PulseCount, delta time.Duration) {
	//log.Printf("pulse %d hb %d", pulse, delta)
	// mobs, scripts, ...

	// pulse zone
	// (zone reset ...)
	if pulse.CheckInterval(mudtime.PulseZone) {
		gs.world.DoZoneActivity()
	}

	// pulse mobs
	// (mobs walk around, initiate attack?)
	if pulse.CheckInterval(mudtime.PulseMobile) {
		gs.world.DoMobileActivity()
	}

	// perform violence
	// do the attacking (players and mobs and everybody)
	if pulse.CheckInterval(mudtime.PulseViolence) {
		gs.world.DoViolence(pulse)
	}

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
			err := gs.handleLogin(msg) // TODO error handling
			if err != nil {
				log.Error().Err(err).Msg("Error from handleLogin")
			}

		case *message.GameMessage_CreatePlayerRequest:
			err := gs.handleCreatePlayer(msg)
			if err != nil {
				log.Error().Err(err).Msg("Error from handleCreatePlayer")
			}

		case *message.GameMessage_DataRequest:
			err := gs.handleDataRequest(msg)
			if err != nil {
				log.Error().Err(err).Msg("Error from handleDataRequest")
			}

		default:
			gs.world.HandleIncomingMessage(msg)
		}
	default:
		// do nothing
	}
	return received
}

func (gs *GameServer) Receive(msg *gameserver.HandlerParameter) {
	gs.incomingBuffer <- msg
}

func (gs *GameServer) Logout(c client.Client, cause string) {
	gm, err := message.NewGameMessage(message.LogoutRequest{Cause: cause})
	if err != nil {
		log.Error().Err(err).Msg("Error creating GameMessage for LogoutRequest")
	} else {
		gs.Receive(gameserver.NewHandlerParameter(c, gm))
	}
}

func (gs *GameServer) handleLogin(msg *gameserver.HandlerParameter) (err error) {
	// is this connection already authenticated?
	// see if we can find an existing player ..
	if msg.Client.Player() != nil {
		// you've already got one
		err = msg.Client.Send(message.LoginResponse{
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

	// todo authentication and stuff - does GRPC have a built in authentication method?

	//playerName := msg.Message.GetLoginRequest().PlayerName
	//playerData, err := db.GetPlayerData(playerName)
	//player := NewClientPlayerFromPlayerData(msg.Message.GetLoginRequest().PlayerName, &playerData, msg.Client)

	// load inventory: have to convert PlayerInventoryData into
	// instances and definitions here, because we need 'the world' to do it.
	/*
		for _, inv := range playerData.Inventory {
			inst, err := gs.world.CreateObjectInstance(inv.ZoneId, inv.DefinitionId, inv.InstanceId)
			if err != nil {
				log.Error().Err(err).Msgf("Error trying to load player %d (%s) inventory instance (%s-%s-%s) -- %s", playerData.Id, playerName, inv.ZoneId, inv.DefinitionId, inv.InstanceId, err)
				clientErr := msg.Client.Send(message.LoginResponse{
					Success:    false,
					ResultCode: "PLAYER_INVENTORY_DATA_ERROR",
				})
				if clientErr != nil {
					log.Error().Err(clientErr).Msg("client error trying to send PLAYER_INVENTORY_DATA_ERROR on login")
				}
				return err
			}
			player.Inventory().Load(inst)
		}
	*/
	// slots - need inventory before we can set slots
	/*
		for _, sd := range playerData.Slots.Slots {
			inst, exists := player.Inventory().GetByInstanceId(sd.InstanceId)
			if !exists {
				log.Error().Msgf("Error trying to load player %d (%s) slot: %d object instance doesn't exist in inventory: %s",
					playerData.Id, playerName, sd.Location, sd.InstanceId)
			} else {
				player.Slots().Set(slot.Location(sd.Location), inst)
			}
		}
	*/

	//msg.Client.SetPlayer(player)
	//msg.Player = player

	// add player to world
	//gs.world.AddPlayer(player)

	//err = player.Send(message.LoginResponse{
	//	Success:    true,
	//	ResultCode: "OK",
	//	PlayerName: player.GetName(),
	//})
	return
}

func (gs *GameServer) handleCreatePlayer(msg *gameserver.HandlerParameter) (err error) {
	if msg.Client.Player() != nil {
		// you've already got one
		err = msg.Client.Send(message.CreatePlayerResponse{
			Success:    false,
			ResultCode: "PLAYER_ALREADY_ATTACHED",
		})
		return
	}
	// TODO fixme
	return
}

// The client is requesting game data: races, class definitions, something like that.
func (gs *GameServer) handleDataRequest(msg *gameserver.HandlerParameter) (err error) {
	resp := message.DataResponse{
		Success:    true,
		ResultCode: "OK",
	}
	resp.DataType = append(resp.DataType, "races")
	// TODO replace all this (or remove it)
	/*
		// get from db
		racejson, err := db.GetRaceDataJson()
		if err != nil {
			log.Error().Err(err).Msg("GetRaceDataJson failed")
			if clientErr := msg.Client.Send(message.DataResponse{
				Success:    false,
				ResultCode: "DATA_ERROR",
			}); clientErr != nil {
				log.Error().Err(clientErr).Msg("handleDataRequest failed to send DB_ERROR for 'races' request")
			}
			return
		}
		resp.Data = append(resp.Data, racejson)
	*/

	// TODO replace all this
	/*
		resp.DataType = append(resp.DataType, "classes")
		classjson, err := db.GetClassDataJson()
		if err != nil {
			log.Error().Err(err).Msg("GetClassDataJson failed")
			if clientErr := msg.Client.Send(message.DataResponse{
				Success:    false,
				ResultCode: "DATA_ERROR",
			}); clientErr != nil {
				log.Error().Err(clientErr).Msg("handleDataRequest failed to send DB_ERROR for 'classes' request")
			}
			return
		}
		resp.Data = append(resp.Data, classjson)
	*/
	if err = msg.Client.Send(resp); err != nil {
		log.Error().Err(err).Msg("handleDataRequest failed to send race data")
	}
	return
}
