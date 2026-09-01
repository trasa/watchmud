package world

import (
	"errors"
	"fmt"
	"iter"
	"log"
	"maps"
	"slices"

	"github.com/google/uuid"
	"github.com/trasa/watchmud-message/direction"
	"github.com/trasa/watchmud/combat"
	"github.com/trasa/watchmud/gameserver"
	"github.com/trasa/watchmud/loader"
	"github.com/trasa/watchmud/mobile"
	"github.com/trasa/watchmud/object"
	"github.com/trasa/watchmud/player"
	"github.com/trasa/watchmud/spaces"
)

// noinspection GoNameStartsWithPackageName
type World struct {
	StartRoom *spaces.Room
	VoidRoom  *spaces.Room
	content   *loader.Content

	// TODO merge playerList and playerRooms similar to MobileRoomMap merges mobList and mobRooms
	playerList  *player.List   // list of players
	playerRooms *PlayerRoomMap // player -> room; room -> players

	mobileRooms *spaces.MobileRoomMap // mobile -> room; room -> mobiles
	handlerMap  map[string]func(message *gameserver.HandlerParameter)

	fightLedger *combat.FightLedger
}

// New creates a brand-new World based on this content
func New(content *loader.Content) (w *World, err error) {
	w = &World{
		content:     content,
		playerList:  player.NewList(),
		playerRooms: NewPlayerRoomMap(),
		mobileRooms: spaces.NewMobileRoomMap(),
		fightLedger: combat.NewFightLedger(),
	}
	w.initializeHandlerMap()
	if err := w.initialLoad(); err != nil {
		return nil, fmt.Errorf("building world: %w", err)
	}
	log.Print("World built.")
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
// Don't send room notifications.
func (w *World) AddPlayer(players ...player.Player) {
	for _, p := range players {
		log.Printf("Adding Player: %s", p.GetName())
		r, exists := w.findRoomByLocation(p.Location())
		if !exists {
			log.Printf("Adding player %s to location %s but it doesn't exist - using start room instead.",
				p.GetName(), p.Location())
			r = w.StartRoom
		}
		w.playerList.Add(p)
		w.playerRooms.Add(p, r)
		r.AddPlayer(p)
	}
}

func (w *World) RemovePlayer(players ...player.Player) {
	for _, p := range players {
		log.Printf("Removing Player: %s", p.GetName())
		w.playerList.Remove(p)
		w.playerRooms.Remove(p)
	}
}

// Player is moving from src room to dest room.
func (w *World) movePlayer(p player.Player, dir direction.Direction, src *spaces.Room, dest *spaces.Room) {
	src.PlayerLeaves(p, dir)
	dest.PlayerEnters(p)
	w.playerRooms.Remove(p)
	w.playerRooms.Add(p, dest)
	p.Location().RoomId = dest.Id
	p.Location().ZoneId = dest.Zone.Id
}

// Player is jumping from the room they are currently in to the destination.
func (w *World) movePlayerMagically(p player.Player, dest *spaces.Room) {
	src := w.getRoomContainingPlayer(p)
	w.movePlayer(p, direction.None, src, dest)
}

// Mobile is moving from src room to dest room.
func (w *World) moveMobile(mob *mobile.Instance, dir direction.Direction, src *spaces.Room, dest *spaces.Room) {
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

func (w *World) getRoomContainingPlayer(p player.Player) *spaces.Room {
	return w.playerRooms.Get(p)
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

func (w *World) findPlayerByName(name string) player.Player {
	return w.playerList.FindByName(name)
}

// Send a message to all players in the world.
func (w *World) SendToAllPlayers(message interface{}) {
	w.playerList.Iter(func(p player.Player) {
		p.Send(message)
	})
}

func (w *World) SendToAllPlayersExcept(exception player.Player, message interface{}) {
	w.playerList.Iter(func(p player.Player) {
		if exception != p {
			p.Send(message)
		}
	})
}

func (w *World) Zones() iter.Seq[*spaces.Zone] {
	return maps.Values(w.content.Zones)
}
func (w *World) Zone(zoneId string) *spaces.Zone {
	return w.content.Zones[zoneId]
}

// CreateObjectInstance builds a new object.Instance for the zoneId, definitionId, and instanceId.
func (w *World) CreateObjectInstance(zoneId string, definitionId string, instanceId uuid.UUID) (*object.Instance, error) {
	z := w.Zone(zoneId)
	d := z.ObjectDefinitions[definitionId]
	return object.NewExistingInstance(instanceId, d)
}
