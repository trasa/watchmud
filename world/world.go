package world

import (
	"github.com/trasa/watchmud/direction"
	"github.com/trasa/watchmud/message"
	"github.com/trasa/watchmud/mobile"
	"github.com/trasa/watchmud/object"
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
		mobs:        make(thing.Map),
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

	// TODO instance ids
	fountainObj := object.NewInstance("fountain", fountainDefn)
	// put the obj in the room
	centralPortalRoom.AddInventory(fountainObj)

	// that's not a knife....wait, yes it is.
	knifeDefn := object.NewDefinition(
		"knife",
		"knife",
		object.WEAPON,
		[]string{},
		"knife",
		"A knife is on the ground.")
	// TODO instance ids
	knifeObj := object.NewInstance("knife", knifeDefn)
	// knife is in room
	centralPortalRoom.AddInventory(knifeObj)

	// a mob - somebody to walk around.
	walkerDefn := mobile.NewDefinition("walker", "walker", []string{}, "The Walker walks.",
		"The walker stands here...for now.",
		true)

	// TODO instance ids
	// TODO put in world as well as room - someone annoying, maybe fix into one method?
	walkerObj := mobile.NewInstance("walker", walkerDefn)
	w.mobs.Add(walkerObj)
	centralPortalRoom.AddMobile(walkerObj)
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

// Walk through all the mob instances that are in this world
// right now and tell them all to do something, if they have
// anything they want to do.
func (w *World) DoMobileActivity() {
	// for each mob in the world
	// wake it up and tell it to do stuff
	// don't limit this to per zone or per room or something
	// remember that mobs can leave the zone they started out in
	// (if programmed to)
	// (or if they really really want to)
	for _, t := range w.mobs {
		mob := t.(*mobile.Instance)
		log.Printf("checking %s", mob.Id())
		if mob.Definition.CanWander {
			log.Printf("can wander! ask about wandering..")
		}
	}
}
