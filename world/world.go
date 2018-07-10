package world

import (
	"github.com/trasa/watchmud/combat"
	"github.com/trasa/watchmud/direction"
	"github.com/trasa/watchmud/gameserver"
	"github.com/trasa/watchmud/mobile"
	"github.com/trasa/watchmud/player"
	"github.com/trasa/watchmud/spaces"
	"log"
)

//noinspection GoNameStartsWithPackageName
type World struct {
	settings  map[string]string
	zones     map[string]*spaces.Zone
	startRoom *spaces.Room
	voidRoom  *spaces.Room

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

// Add Player(s) to the world putting them in the start room,
// Don't send room notifications.
func (w *World) AddPlayer(players ...player.Player) {
	for _, p := range players {
		log.Printf("Adding Player: %s", p.GetName())
		w.playerList.Add(p)
		w.playerRooms.Add(p, w.startRoom)
		w.startRoom.AddPlayer(p)
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
}

// Mobile is moving from src room to dest room.
func (w *World) moveMobile(mob *mobile.Instance, dir direction.Direction, src *spaces.Room, dest *spaces.Room) {
	src.MobileLeaves(mob, dir)
	dest.MobileEnters(mob)
	w.mobileRooms.Remove(mob)
	w.mobileRooms.Add(mob, dest)
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
