package world

import (
	"github.com/trasa/watchmud/direction"
	"github.com/trasa/watchmud/message"
	"github.com/trasa/watchmud/object"
	"github.com/trasa/watchmud/player"
	"log"
)

//noinspection GoNameStartsWithPackageName
type World struct {
	zones       map[string]*Zone
	startRoom   *Room
	voidRoom    *Room
	playerList  *player.List
	playerRooms *PlayerRoomMap
	handlerMap  map[string]func(message *message.IncomingMessage)
}

// Constructor for World
func New() *World {
	// Build a very boring world.
	w := World{
		zones:       make(map[string]*Zone),
		playerList:  player.NewList(),
		playerRooms: NewPlayerRoomMap(),
	}
	w.initializeHandlerMap()

	sampleZone := Zone{
		Id:    "sample",
		Name:  "Sample Zone",
		Rooms: make(map[string]*Room),
	}
	w.zones[sampleZone.Id] = &sampleZone

	/*
		  north -- northeast
		   |         |
		central -- east

	*/
	// central room (Start)
	centralPortalRoom := NewRoom(&sampleZone, "start", "Central Portal", "It's a boring room, with boring stuff in it.")
	sampleZone.Rooms[centralPortalRoom.Id] = centralPortalRoom
	w.startRoom = centralPortalRoom

	// north room
	northRoom := NewRoom(&sampleZone, "northRoom", "North Room", "This room is north of the start.")
	sampleZone.Rooms[northRoom.Id] = northRoom

	// northeast
	northeastRoom := NewRoom(&sampleZone, "northeastRoom", "North East Room", "It's north, and also East.")
	sampleZone.Rooms[northeastRoom.Id] = northeastRoom

	// east
	eastRoom := NewRoom(&sampleZone, "eastRoom", "East Room", "This room is east of the start.")
	sampleZone.Rooms[eastRoom.Id] = eastRoom

	// central <-> north
	centralPortalRoom.North = northRoom
	northRoom.South = centralPortalRoom

	// central <-> east
	centralPortalRoom.East = eastRoom
	eastRoom.West = centralPortalRoom

	// north <-> northeast
	northRoom.East = northeastRoom
	northeastRoom.West = northRoom

	// east <-> northeast
	eastRoom.North = northeastRoom
	northeastRoom.South = eastRoom

	// The VOID. When you're not really in a room.
	w.voidRoom = NewRoom(nil, "void", "The Void", "You see nothing but endless void.")

	// lets put "something" in the central portal room
	fountainDefn := object.NewDefinition(
		"fountain",
		"fountain",
		object.OTHER,
		[]string{"fount"},
		"fountain",
		"A fountain bubbles quietly.")

	fountainObj := object.NewInstance(fountainDefn)
	// put the obj in the room
	centralPortalRoom.AddObject(fountainObj)

	// that's not a knife....wait, yes it is.
	knifeDefn := object.NewDefinition(
		"knife",
		"knife",
		object.WEAPON,
		[]string{},
		"knife",
		"A knife is on the ground.")
	knifeObj := object.NewInstance(knifeDefn)
	// knife is in room
	centralPortalRoom.AddObject(knifeObj)

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
		w.startRoom.Add(p)
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
	src.Leave(p, dir)
	dest.Enter(p)
	w.playerRooms.Remove(p)
	w.playerRooms.Add(p, dest)
}

func (w *World) getRoomContainingPlayer(p player.Player) *Room {
	return w.playerRooms.Get(p)
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
