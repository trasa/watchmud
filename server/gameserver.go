package server

import (
	"context"
	"errors"
	"fmt"
	"time"
	"uuid"

	"github.com/rs/zerolog/log"
	"github.com/trasa/watchmud-message"
	"github.com/trasa/watchmud/gameserver"
	"github.com/trasa/watchmud/mudtime"
	"github.com/trasa/watchmud/player"
	"github.com/trasa/watchmud/rules"
	"github.com/trasa/watchmud/world"
)

type GameServer struct {
	incomingBuffer chan *gameserver.HandlerParameter
	world          *world.World
	catalog        *rules.Catalog
	tickInterval   time.Duration
	store          player.Store
}

func New(w *world.World, c *rules.Catalog, s player.Store) *GameServer {
	const bufferSize = 64
	return &GameServer{
		incomingBuffer: make(chan *gameserver.HandlerParameter, bufferSize),
		world:          w,
		catalog:        c,
		store:          s,
	}
}

// Run the game server, obviously.
func (gs *GameServer) Run(ctx context.Context) error {
	ticker := time.NewTicker(mudtime.PulseInterval)
	defer ticker.Stop()

	last := time.Now()
	var pulse mudtime.PulseCount

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case msg := <-gs.incomingBuffer:
			if err := gs.dispatch(msg); err != nil {
				// don't return error, we're not halting the server
				log.Error().Err(err).Msg("error dispatching message")
			}
		case <-ticker.C:
			now := time.Now()
			delta := now.Sub(last)
			last = now
			pulse++
			gs.heartbeat(pulse, delta) // zone/mob/violence pulses
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
}

// dispatch a message to its handler.
func (gs *GameServer) dispatch(msg *gameserver.HandlerParameter) error {
	switch msg.Message.Inner.(type) {
	case *message.GameMessage_LoginRequest:
		if err := gs.handleLogin(msg); err != nil {
			return err
		}
	case *message.GameMessage_CreatePlayerRequest:
		if err := gs.handleCreatePlayer(msg); err != nil {
			return err
		}
	case *message.GameMessage_DataRequest:
		if err := gs.handleDataRequest(msg); err != nil {
			return err
		}
	default:
		if err := gs.world.HandleIncomingMessage(msg); err != nil {
			return err
		}
	}
	return nil
}

func (gs *GameServer) Receive(msg *gameserver.HandlerParameter) {
	gs.incomingBuffer <- msg
}

func (gs *GameServer) Logout(c gameserver.Conn, cause string) {
	gm, err := message.NewGameMessage(message.LogoutRequest{Cause: cause})
	if err != nil {
		log.Error().Err(err).Msg("Error creating GameMessage for LogoutRequest")
	} else {
		gs.Receive(gameserver.NewHandlerParameter(c, gm))
	}
}

func (gs *GameServer) handleLogin(msg *gameserver.HandlerParameter) error {
	// is this connection already authenticated?
	// see if we can find an existing player.
	if msg.Client.Player() != nil {
		// you've already got one - this is an error in our connection logic
		return errors.New("player already attached to client")
	}

	// what if player is logged in on a different client?
	/*
		if p := FindPlayerByClient(message.Client); p != nil {
			// TODO: kick the old user and proceed with the new
			// for now, fail the login
			return errors.New("player already logged in")
		}*/

	// TODO authentication and stuff...
	playerName := msg.Message.GetLoginRequest().PlayerName
	rec, found, err := gs.store.Load(playerName)
	if err != nil {
		// store error - problem with the store, return an error
		return err
	}
	if !found {
		// not an error - could represent a new player (player creation)
		log.Info().Str("playerName", playerName).Msg("playerName not found in store")
		loginResponse := message.LoginResponse{
			Success:    false,
			ResultCode: "PLAYER_LOGIN_FAILED",
		}
		msg.Client.Send(loginResponse)
		return nil
	}

	// create the player
	p, err := player.FromRecord(rec, msg.Client, gs.catalog, gs.world)
	if err != nil {
		return err
	}
	msg.Player = p
	msg.Client.SetPlayer(p)

	// add player to world
	gs.world.AddPlayer(p)

	p.Send(message.LoginResponse{
		Success:    true,
		ResultCode: "OK",
		PlayerName: p.Name(),
	})
	return nil
}

func (gs *GameServer) handleCreatePlayer(msg *gameserver.HandlerParameter) error {
	if msg.Client.Player() != nil {
		// you've already got one
		// this is a programming bug (login state machine), so report the error
		return fmt.Errorf("player %s already attached to client", msg.Client.Player().Name())
	}
	req := msg.Message.GetCreatePlayerRequest()
	playerName := req.PlayerName
	// message sends an int but this is a string now, so default all to human
	lineage := gs.catalog.Lineages["human"]
	// message sends an int but this is a string now, so default all to fighter
	class := gs.catalog.Classes["fighter"]
	p := player.New(
		uuid.New(),
		playerName,
		msg.Client,
		lineage,
		class,
		rules.StandardAbilities(class.AbilityPreference),
	)

	// TODO need to set the location first (AddPlayer always puts the player in the start room, for now)

	if err := gs.store.Save(p.Record()); err != nil {
		return fmt.Errorf("handleCreatePlayer: %v", err)
	}

	msg.Client.SetPlayer(p)
	msg.Player = p

	gs.world.AddPlayer(p)

	p.Send(message.CreatePlayerResponse{
		Success:    true,
		ResultCode: "OK",
		PlayerName: p.Name(),
	})
	return nil
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
	msg.Client.Send(resp)
	return
}
