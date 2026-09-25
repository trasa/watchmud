package world

import (
	"errors"
	"fmt"
	"iter"
	"maps"
	"slices"

	"github.com/rs/zerolog/log"
	"github.com/watchmud/watchmud/combat"
	"github.com/watchmud/watchmud/event"
	"github.com/watchmud/watchmud/loader"
	"github.com/watchmud/watchmud/mobile"
	"github.com/watchmud/watchmud/object"
	"github.com/watchmud/watchmud/player"
	"github.com/watchmud/watchmud/rules"
	"github.com/watchmud/watchmud/spaces"
)

// noinspection GoNameStartsWithPackageName
type World struct {
	StartRoom *spaces.Room
	VoidRoom  *spaces.Room
	DeathRoom *spaces.Room // where a player who dies wakes up
	content   *loader.Content

	roller rules.Roller
	store  player.Store

	playerList   *player.List   // list of players in world
	playerToRoom *playerRoomMap // map of player to the room they are in. rooms own their player list.

	mobileRooms *spaces.MobileRoomMap // mobile -> room; room -> mobiles

	fightLedger *combat.FightLedger

	reservedNames map[string]bool // by player.NameKey; see IsReservedName
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
	w.reservedNames = reservedNames(c)
	if err := w.initialLoad(); err != nil {
		return nil, fmt.Errorf("building world: %w", err)
	}
	log.Info().Msg("World built.")
	return w, nil
}

func (w *World) initialLoad() (err error) {
	content := w.content
	settings := content.Settings
	if w.StartRoom, err = content.Room(settings.Start.ZoneId, settings.Start.RoomId); err != nil {
		return fmt.Errorf("start room: %w", err)
	}
	if w.VoidRoom, err = content.Room(settings.Void.ZoneId, settings.Void.RoomId); err != nil {
		return fmt.Errorf("void room: %w", err)
	}
	if w.DeathRoom, err = content.Room(settings.PlayerDeath.ZoneId, settings.PlayerDeath.RoomId); err != nil {
		return fmt.Errorf("player-death room: %w", err)
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

// AddPlayer or players to the world, in the start room. New characters come
// in this way; returning ones go through ReturnPlayer.
func (w *World) AddPlayer(players ...*player.Player) {
	for _, p := range players {
		w.addPlayerTo(p, w.StartRoom)
	}
}

// ReturnPlayer puts a returning player back in the room their record says
// they were last in, or the start room if it doesn't say or that room no
// longer exists. A content edit that removes a room must not strand anybody,
// the same as the missing object definitions player.FromRecord forgives.
func (w *World) ReturnPlayer(p *player.Player, zoneId, roomId string) {
	r, found := w.findRoomById(zoneId, roomId)
	if !found {
		if zoneId != "" || roomId != "" {
			p.Log().Warn().Msgf("last room %s.%s not found, using the start room", zoneId, roomId)
		}
		r = w.StartRoom
	}
	w.addPlayerTo(p, r)
}

// Arrive finishes a login: the room hears who just appeared in it, and the
// player is shown where they are -- which, now that a returning player comes
// back to wherever they left, isn't necessarily where they expect.
//
// Separate from AddPlayer and ReturnPlayer so the transport's login events go
// out first; the player's own description has to land after the login
// conversation has ended, not in the middle of it.
func (w *World) Arrive(p *player.Player) {
	r := w.getPlayerRoom(p)
	r.SendExcept(p, event.EnteredGame{Actor: p.Name()})
	p.Send(r.DescriptionExcept(p))
}

func (w *World) addPlayerTo(p *player.Player, r *spaces.Room) {
	p.Log().Debug().Str("room", r.Location().String()).Msg("Adding player to world")
	w.playerList.Add(p)
	r.AddPlayer(p)
	w.playerToRoom.Update(p, r)
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

// AddMobile adds a mobile instance to a room in the world.
func (w *World) AddMobile(mob *mobile.Instance, targetRoom *spaces.Room) {
	w.mobileRooms.Add(mob, targetRoom)
}

// RemoveMobile removes a mobile instance from the world.
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

// IsPlaying says whether a character by that name is in the world right now.
func (w *World) IsPlaying(name string) bool {
	return w.findPlayerByName(name) != nil
}

func (w *World) findPlayerByName(name string) *player.Player {
	return w.playerList.FindByName(name)
}

// Players is everyone in the world.
func (w *World) Players() iter.Seq[*player.Player] {
	return w.playerList.All()
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
