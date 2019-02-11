package world

import (
	"github.com/satori/go.uuid"
	"github.com/trasa/watchmud-message/direction"
	"github.com/trasa/watchmud/combat"
	"github.com/trasa/watchmud/gameserver"
	"github.com/trasa/watchmud/mobile"
	"github.com/trasa/watchmud/object"
	"github.com/trasa/watchmud/player"
	"github.com/trasa/watchmud/spaces"
	"log"
)

//noinspection GoNameStartsWithPackageName
type World struct {
	settings  map[string]string
	zones     map[string]*spaces.Zone
	StartRoom *spaces.Room
	VoidRoom  *spaces.Room

	// TODO merge playerList and playerRooms similar to MobileRoomMap merges mobList and mobRooms
	playerList  *player.List   // list of players
	playerRooms *PlayerRoomMap // player -> room; room -> players

	mobileRooms *spaces.MobileRoomMap // mobile -> room; room -> mobiles
	handlerMap  map[string]func(message *gameserver.HandlerParameter)

	fightLedger *combat.FightLedger
}

// Constructor for World
func New(worldFilesDirectory string) (w *World, err error) {
	// Build a very boring world.
	w = &World{
		zones:       make(map[string]*spaces.Zone),
		playerList:  player.NewList(),
		playerRooms: NewPlayerRoomMap(),
		mobileRooms: spaces.NewMobileRoomMap(),
		fightLedger: combat.NewFightLedger(),
	}
	w.initializeHandlerMap()
	err = w.initialLoad(worldFilesDirectory)
	log.Print("World built.")
	if err != nil {
		log.Printf("world.New: There were errors: %v", err)
	}
	return
}

// Add Player(s) to the world putting them in the correct room they were
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

// Mobile is moving from src room to dest room.
func (w *World) moveMobile(mob *mobile.Instance, dir direction.Direction, src *spaces.Room, dest *spaces.Room) {
	src.MobileLeaves(mob, dir)
	dest.MobileEnters(mob)
	w.mobileRooms.Remove(mob)
	w.mobileRooms.Add(mob, dest)
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

// find room by zone id and room id.
// return nil if not found
func (w *World) findRoomById(zoneId string, roomId string) (*spaces.Room, bool) {
	if z, zoneExists := w.zones[zoneId]; zoneExists {
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

// Send a message to all players in the world except for one.
func (w *World) SendToAllPlayersExcept(exception player.Player, message interface{}) {
	w.playerList.Iter(func(p player.Player) {
		if exception != p {
			p.Send(message)
		}
	})
}

func (w *World) GetZone(zoneId string) *spaces.Zone {
	return w.zones[zoneId]
}

// Create an object.Instance for this zone, definition, and instance ID
func (w *World) CreateObjectInstance(zoneId string, definitionId string, instanceId uuid.UUID) (*object.Instance, error) {
	z := w.GetZone(zoneId)
	defn := z.ObjectDefinitions[definitionId]
	return object.NewExistingInstance(instanceId, defn)
}
