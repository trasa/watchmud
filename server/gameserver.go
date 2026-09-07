package server

import (
	"context"
	"errors"
	"time"
	"uuid"

	"github.com/rs/zerolog/log"
	"github.com/trasa/watchmud-message"
	"github.com/trasa/watchmud/client"
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
	return &GameServer{
		incomingBuffer: make(chan *gameserver.HandlerParameter),
		world:          w,
		catalog:        c,
		store:          s,
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

func (gs *GameServer) handleLogin(msg *gameserver.HandlerParameter) error {
	// is this connection already authenticated?
	// see if we can find an existing player ..
	if msg.Client.Player() != nil {
		// you've already got one
		// TODO error handling on send
		msg.Client.Send(message.LoginResponse{
			Success:    false,
			ResultCode: "PLAYER_ALREADY_ATTACHED",
		})
		return errors.New("player already attached to client")
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

	// TODO authentication and stuff...
	playerName := msg.Message.GetLoginRequest().PlayerName
	rec, found, loadErr := gs.store.Load(playerName)
	if loadErr != nil {
		// store error
		log.Error().Err(loadErr).Str("playerName", playerName).Msgf("Error loading player %s from store", playerName)
		loginResponse := message.LoginResponse{
			Success:    false,
			ResultCode: "PLAYER_STORE_ERROR",
		}
		if err := msg.Client.Send(loginResponse); err != nil {
			log.Error().Err(err).Msg("client error trying to send PLAYER_STORE_ERROR on login")
		}
		return loadErr
	}
	if !found {
		log.Warn().Str("playerName", playerName).Msgf("playerName %s not found in store", playerName)
		loginResponse := message.LoginResponse{
			Success:    false,
			ResultCode: "PLAYER_LOGIN_FAILED",
		}
		if err := msg.Client.Send(loginResponse); err != nil {
			log.Error().Err(err).Msg("client error trying to send PLAYER_LOGIN_FAILED on login")
			// TODO deal with send error
		}
		// TODO should this return an error?
		return errors.New("player not found in store")
	}

	p, err := player.FromRecord(rec, msg.Client, gs.catalog, gs.world)
	if err != nil {
		log.Error().Err(err).Msg("Error creating player from record")
	}
	msg.Player = p

	// add player to world
	gs.world.AddPlayer(p)

	if err := p.Send(message.LoginResponse{
		Success:    true,
		ResultCode: "OK",
		PlayerName: p.Name,
	}); err != nil {
		log.Error().Err(err).Msg("Error sending LoginResponse")
		return err
	}
	return nil
}

func (gs *GameServer) handleCreatePlayer(msg *gameserver.HandlerParameter) error {
	if msg.Client.Player() != nil {
		// you've already got one
		// TODO error handling for send (there's not really much we can honestly do ... so don't return an error?)
		// Or better, consolidate all the error handling back in what calls this
		msg.Client.Send(message.CreatePlayerResponse{
			Success:    false,
			ResultCode: "PLAYER_ALREADY_ATTACHED",
		})
		return errors.New("player already attached")
	}
	req := msg.Message.GetCreatePlayerRequest()
	playerName := req.PlayerName
	// message sends an int but this is a string now, so default all to human
	lineage := gs.catalog.Lineages["human"]
	// message sends an int but this is a string now, so default all to fighter
	class := gs.catalog.Classes["fighter"]
	// TODO need check for name uniqueness, and shouldn't use it as the ID... need uuid support
	// NOTE this doesn't set the location that's defered to the world I think...?
	p := player.New(
		uuid.New(),
		playerName,
		msg.Client,
		lineage,
		class,
		rules.StandardAbilities(class.AbilityPreference),
	)
	// TODO how do we set Client.Player() now that we've changed things
	// msg.Client.Player = p
	msg.Player = p

	// TODO need to set the location first (AddPlayer always puts the player in the start room, for now)

	if err := gs.store.Save(p.Record()); err != nil {
		log.Error().Err(err).Msgf("Error trying to save player record for %s", playerName)
		// TODO send error to client
		return errors.New("error saving player record")
	}
	gs.world.AddPlayer(p)

	err := p.Send(message.CreatePlayerResponse{
		Success:    true,
		ResultCode: "OK",
		PlayerName: p.Name,
	})
	return err
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
