package server

import (
	"context"
	"errors"
	"fmt"
	"time"
	"uuid"

	"github.com/rs/zerolog/log"
	"github.com/trasa/watchmud/command"
	"github.com/trasa/watchmud/event"
	"github.com/trasa/watchmud/gameserver"
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
	log.Info().Msg("starting game server Run loop")

	ticker := time.NewTicker(rules.PulseInterval)
	defer ticker.Stop()

	last := time.Now()
	var pulse rules.PulseCount

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case msg := <-gs.incomingBuffer:
			if err := gs.dispatch(msg); err != nil {
				// don't return error, we're not halting the server
				log.Error().Err(err).Msg("error dispatching message")
			}
			gs.prompt()
		case <-ticker.C:
			now := time.Now()
			delta := now.Sub(last)
			last = now
			pulse++
			gs.heartbeat(pulse, delta) // zone/mob/violence pulses
			gs.prompt()
		}
	}
}

// prompt tells every player the game is waiting on them again. Everyone, not
// just whoever sent the command: a say or a combat round lands on bystanders
// too, and they need their prompt back as much as the speaker does. The
// transport drops the prompts that nothing was said in front of, so this
// costs a channel send per player and prints nothing new to the rest.
//
// A player who has just logged in is in the list by now, so this is also
// their first prompt; one whose login failed is not, and gets none.
func (gs *GameServer) prompt() {
	for p := range gs.world.Players() {
		p.Send(event.Prompt{CurrentHealth: p.CurrentHealth(), MaxHealth: p.MaxHealth()})
	}
}

// Heartbeat runner of the game. Use pulse to determine intervals
// between things (ex. reset zones every 15 minutes...)
// delta is the amount of time since the last heartbeat was run.
func (gs *GameServer) heartbeat(pulse rules.PulseCount, delta time.Duration) {
	//log.Printf("pulse %d hb %d", pulse, delta)
	// mobs, scripts, ...

	// pulse zone
	// (zone reset ...)
	zonePulse := gs.catalog.MudTime.Zone
	if pulse.CheckInterval(zonePulse) {
		gs.world.DoZoneActivity()
	}

	// pulse mobs
	// (mobs walk around, initiate attack?)
	mobPulse := gs.catalog.MudTime.Mobile
	if pulse.CheckInterval(mobPulse) {
		gs.world.DoMobileActivity()
	}

	// perform violence
	// do the attacking (players and mobs and everybody)
	violencePulse := gs.catalog.MudTime.Violence
	if pulse.CheckInterval(violencePulse) {
		gs.world.DoViolence(pulse)
	}

	// mud-hour ("player tick")
	// affect weather, regen ..

	// saving player data
	// if pulse.CheckInterval(mudtime.)
}

// dispatch a message to its handler.
//
// The pre-world cases -- login, create player -- are handled here because
// they need the store and player construction. Everything else falls through
// to the world.
func (gs *GameServer) dispatch(msg *gameserver.HandlerParameter) error {
	switch cmd := msg.Command.(type) {
	case command.Login:
		return gs.handleLogin(msg, cmd)
	case command.CreatePlayer:
		return gs.handleCreatePlayer(msg, cmd)
	default:
		return gs.world.HandleIncomingMessage(msg)
	}
}

func (gs *GameServer) Receive(msg *gameserver.HandlerParameter) {
	gs.incomingBuffer <- msg
}

func (gs *GameServer) Logout(c gameserver.Conn, cause string) {
	gs.Receive(gameserver.NewHandlerParameter(c, command.Logout{Cause: cause}))
}

func (gs *GameServer) handleLogin(msg *gameserver.HandlerParameter, cmd command.Login) error {
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
	playerName := cmd.Name
	rec, found, err := gs.store.Load(playerName)
	if err != nil {
		// store error - problem with the store, return an error
		return err
	}
	if !found {
		// not an error - could represent a new player (player creation)
		log.Info().Str("playerName", playerName).Msg("playerName not found in store")
		msg.Client.Send(event.LoginFailed{Reason: event.NoSuchPlayer})
		return nil
	}

	// create the player
	p, err := player.FromRecord(rec, msg.Client, gs.catalog, gs.world)
	if err != nil {
		return fmt.Errorf("handleLogin %s: %w", playerName, err)
	}
	msg.Player = p
	msg.Client.SetPlayer(p)

	// back where they left off
	gs.world.ReturnPlayer(p, rec.LastZoneId, rec.LastRoomId)

	p.Send(event.LoggedIn{Name: p.Name()})
	gs.world.Arrive(p)
	return nil
}

func (gs *GameServer) handleCreatePlayer(msg *gameserver.HandlerParameter, cmd command.CreatePlayer) error {
	if msg.Client.Player() != nil {
		// you've already got one
		// this is a programming bug (login state machine), so report the error
		return fmt.Errorf("player %s already attached to client", msg.Client.Player().Name())
	}
	playerName := cmd.Name

	// The lineage is the only choice creation makes, and it is cosmetic. An
	// id the catalog doesn't know means the transport offered something stale
	// -- worth a log line, not worth refusing to make the character.
	lineage, found := gs.catalog.Lineages[cmd.Lineage]
	if !found {
		if cmd.Lineage != "" {
			log.Warn().Str("playerName", playerName).Msgf("unknown lineage %q at creation, using the default", cmd.Lineage)
		}
		lineage = gs.catalog.DefaultLineage()
	}
	if lineage == nil {
		return errors.New("handleCreatePlayer: no lineages defined in the catalog")
	}

	p := player.New(
		uuid.New(),
		playerName,
		msg.Client,
		lineage,
		gs.catalog,
	)

	// The kit goes on before the save, so the first thing written for this
	// character already has it: a crash between here and the first command
	// leaves them dressed rather than naked.
	player.GiveStartingGear(p, gs.catalog.StartingGear, gs.world)

	// TODO need to set the location first (AddPlayer always puts the player in the start room, for now)

	if err := gs.store.Save(p.Record()); err != nil {
		return fmt.Errorf("handleCreatePlayer: %v", err)
	}

	msg.Client.SetPlayer(p)
	msg.Player = p

	gs.world.AddPlayer(p)

	p.Send(event.PlayerCreated{Name: p.Name()})
	gs.world.Arrive(p)
	return nil
}
