package world

import (
	"errors"
	"fmt"
	"iter"
	"maps"
	"slices"

	"github.com/rs/zerolog/log"
	"github.com/trasa/watchmud/combat"
	"github.com/trasa/watchmud/loader"
	"github.com/trasa/watchmud/mobile"
	"github.com/trasa/watchmud/object"
	"github.com/trasa/watchmud/player"
	"github.com/trasa/watchmud/rules"
	"github.com/trasa/watchmud/spaces"
)

// noinspection GoNameStartsWithPackageName
type World struct {
	StartRoom *spaces.Room
	VoidRoom  *spaces.Room
	content   *loader.Content

	roller rules.Roller
	store  player.Store

	playerList   *player.List   // list of players in world
	playerToRoom *playerRoomMap // map of player to the room they are in. rooms own their player list.

	mobileRooms *spaces.MobileRoomMap // mobile -> room; room -> mobiles

	fightLedger *combat.FightLedger
}

// New creates a brand-new World based on this content
func New(c *loader.Content, s player.Store, roller rules.Roller) (w *World, err error) {
	w = &World{
		content:      c,
		playerList:   player.NewList(),
		playerToRoom: newPlayerRoomMap(),
		mobileRooms:  spaces.NewMobileRoomMap(),
		fightLedger:  combat.NewFightLedger(),
		roller:       roller,
		store:        s,
	}
	if err := w.initialLoad(); err != nil {
		return nil, fmt.Errorf("building world: %w", err)
	}
	log.Info().Msg("World built.")
	return w, nil
}

func (w *World) initialLoad() (err error) {
	content := w.content
	if w.StartRoom, err = content.Room(content.Settings.StartZone, content.Settings.StartRoom); err != nil {
		return fmt.Errorf("start room: %w", err)
	}
	if w.VoidRoom, err = content.Room(content.Settings.VoidZone, content.Settings.VoidRoom); err != nil {
		return fmt.Errorf("void room: %w", err)
	}

	// Process the zone commands that say which
	// mob and object instances to create and where. Distinct from building the
	// world, since this recurs throughout runtime.
	for _, zoneId := range slices.Sorted(maps.Keys(content.Zones)) {
		if errs := content.Zones[zoneId].Reset(w.mobileRooms); len(errs) > 0 {
			return fmt.Errorf("initial reset of zone %s: %w", zoneId, errors.Join(errs...))
		}
	}
	return nil
}

// AddPlayer or players to the world putting them in the correct room they were
// in last time, or the start room if we can't figure that out.
func (w *World) AddPlayer(players ...*player.Player) {
	for _, p := range players {
		p.Log().Debug().Msg("Adding player to world")
		// list of known players
		w.playerList.Add(p)

		// TODO need support for location
		// so for now, this *always* adds to the start room.
		r := w.StartRoom
		r.AddPlayer(p)
		w.playerToRoom.Update(p, r)
	}
}

func (w *World) RemovePlayer(players ...*player.Player) {
	for _, p := range players {
		p.Log().Debug().Msg("Removing player")
		w.fightLedger.EndAllFightsWith(p.Id())
		// in a room? remove it.
		r := w.playerToRoom.Get(p)
		if r != nil {
			r.RemovePlayer(p)
		}
		w.playerList.Remove(p)
		w.playerToRoom.Remove(p)
	}
}

// Player is moving from src room to dest room.
func (w *World) movePlayer(p *player.Player, dir rules.Direction, src *spaces.Room, dest *spaces.Room) {
	src.PlayerLeaves(p, dir)
	dest.PlayerEnters(p)
	w.playerToRoom.Update(p, dest)
}

// Player is jumping from the room they are currently in to the destination.
func (w *World) movePlayerMagically(p *player.Player, dest *spaces.Room) {
	src := w.playerToRoom.Get(p)
	if src == nil {
		src = w.VoidRoom
	}
	w.movePlayer(p, rules.DirectionNone, src, dest)
}

// Mobile is moving from src room to dest room.
func (w *World) moveMobile(mob *mobile.Instance, dir rules.Direction, src *spaces.Room, dest *spaces.Room) {
	src.MobileLeaves(mob, dir)
	dest.MobileEnters(mob)
	w.mobileRooms.Remove(mob)
	w.mobileRooms.Add(mob, dest)
}

// add a mobile instance to the world
func (w *World) AddMobile(mob *mobile.Instance, targetRoom *spaces.Room) {
	w.mobileRooms.Add(mob, targetRoom)
}

// remove the mobile instance from the world entirely
func (w *World) removeMobile(mob *mobile.Instance) {
	w.mobileRooms.Remove(mob)
}

// roleName resolves a player's equipment weights to a role's display name,
// empty when the gear adds up to no role at all. A role is never stored, so
// every caller that wants one derives it here -- see CLAUDE.md, "Lineage and
// Role".
func (w *World) roleName(weights map[string]int) string {
	if r := w.content.Catalog.RoleFor(weights); r != nil {
		return r.Name
	}
	return ""
}

// getPlayerRoom returns the room a player is in, or VoidRoom if we can't figure that out.
// Does not return nil.
func (w *World) getPlayerRoom(p *player.Player) *spaces.Room {
	r := w.playerToRoom.Get(p)
	if r == nil {
		p.Log().Warn().Msg("player not in a room!")
		return w.VoidRoom
	}
	return r
}

func (w *World) getRoomContainingMobile(mob *mobile.Instance) *spaces.Room {
	return w.mobileRooms.GetRoomForMobile(mob)
}

// Find room by zone id and room id.
func (w *World) findRoomById(zoneId string, roomId string) (*spaces.Room, bool) {
	if z, zoneExists := w.content.Zones[zoneId]; zoneExists {
		if r, roomExists := z.Rooms[roomId]; roomExists {
			return r, true
		}
	}
	return nil, false
}

func (w *World) findRoomByLocation(loc *player.Location) (*spaces.Room, bool) {
	if loc == nil {
		return nil, false
	}
	return w.findRoomById(loc.ZoneId, loc.RoomId)
}

func (w *World) findPlayerByName(name string) *player.Player {
	return w.playerList.FindByName(name)
}

// SendToAllPlayers send a message to all players in the world.
func (w *World) SendToAllPlayers(message interface{}) {
	for p := range w.playerList.All() {
		p.Send(message)
	}
}

// SendToAllPlayersExcept send a message to all players in the world except the exception player.
func (w *World) SendToAllPlayersExcept(exception *player.Player, message interface{}) {
	for p := range w.playerList.AllExcept(exception) {
		p.Send(message)
	}
}

func (w *World) Zones() iter.Seq[*spaces.Zone] {
	return maps.Values(w.content.Zones)
}
func (w *World) Zone(zoneId string) *spaces.Zone {
	return w.content.Zones[zoneId]
}

// ObjectDefinition looks up by zoneId and definitionId, implementing player.DefinitionSource
func (w *World) ObjectDefinition(zoneId, definitionId string) (*object.Definition, bool) {
	z := w.Zone(zoneId)
	if z == nil {
		return nil, false
	}
	d, ok := z.ObjectDefinitions[definitionId]
	if !ok {
		return nil, false
	}
	return d, true
}
