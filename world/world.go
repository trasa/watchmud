package world

import (
	"github.com/trasa/watchmud/direction"
	"github.com/trasa/watchmud/message"
	"github.com/trasa/watchmud/mobile"
	"github.com/trasa/watchmud/player"
	"github.com/trasa/watchmud/thing"
	"log"
)

//noinspection GoNameStartsWithPackageName
type World struct {
	zones       map[string]*Zone
	startRoom   *Room
	voidRoom    *Room
	playerList  *player.List
	playerRooms *PlayerRoomMap
	mobileRooms *MobileRoomMap
	handlerMap  map[string]func(message *message.IncomingMessage)
	mobs        thing.Map
}

// Constructor for World
func New() *World {
	// Build a very boring world.
	w := World{
		zones:       make(map[string]*Zone),
		playerList:  player.NewList(),
		playerRooms: NewPlayerRoomMap(),
		mobileRooms: NewMobileRoomMap(),
		mobs:        make(thing.Map),
	}
	w.initializeHandlerMap()
	w.initialLoad()
	log.Print("World built.")
	return &w
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
func (w *World) movePlayer(p player.Player, dir direction.Direction, src *Room, dest *Room) {
	src.PlayerLeaves(p, dir)
	dest.PlayerEnters(p)
	w.playerRooms.Remove(p)
	w.playerRooms.Add(p, dest)
}

// Add mobile(s) to the world putting them in the room indicated,
// Don't send room notifications.
func (w *World) AddMobiles(destRoom *Room, mobs ...*mobile.Instance) {
	for _, mob := range mobs {
		log.Printf("Adding Mobile: %s", mob.InstanceId)
		w.mobs.Add(mob)
		w.mobileRooms.Add(mob, destRoom)
		destRoom.AddMobile(mob)
	}
}

// Mobile is moving from src room to dest room.
func (w *World) moveMobile(mob *mobile.Instance, dir direction.Direction, src *Room, dest *Room) {
	src.MobileLeaves(mob, dir)
	dest.MobileEnters(mob)
	w.mobileRooms.Remove(mob)
	w.mobileRooms.Add(mob, dest)
}

func (w *World) getRoomContainingPlayer(p player.Player) *Room {
	return w.playerRooms.Get(p)
}

func (w *World) getRoomContainingMobile(mob *mobile.Instance) *Room {
	return w.mobileRooms.Get(mob)
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
