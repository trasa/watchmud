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
}

// Constructor for World
func New() *World {
	// Build a very boring world.
	w := World{
		zones:       make(map[string]*Zone),
		playerList:  player.NewList(),
		playerRooms: NewPlayerRoomMap(),
	}
	sampleZone := Zone{
		Id:    "sample",
		Name:  "Sample Zone",
		Rooms: make(map[string]*Room),
	}
	w.zones[sampleZone.Id] = &sampleZone

	centralPortalRoom := NewRoom(&sampleZone, "start", "Central Portal", "It's a boring room, with boring stuff in it.")
	sampleZone.Rooms[centralPortalRoom.Id] = centralPortalRoom
	w.startRoom = centralPortalRoom

	// north room
	northRoom := NewRoom(&sampleZone, "northRoom", "North Room", "This room is north of the start.")
	sampleZone.Rooms[northRoom.Id] = northRoom

	// north room and central portal connect to each other
	centralPortalRoom.North = northRoom
	northRoom.South = centralPortalRoom

	// The VOID. When you're not really in a room.
	w.voidRoom = NewRoom(nil, "void", "The Void", "You see nothing but endless void.")

	// lets put "something" in the central portal room
	fountainDefn := object.NewDefinition(
		"fountain",
		"fountain",
		object.OTHER,
		[]string{"fount"},
		"fountain",
		"A fountain burbles quietly.")

	fountainObj := object.NewInstance(fountainDefn)
	// put the obj in the room
	centralPortalRoom.AddObject(fountainObj)

	log.Print("World built.")
	return &w
}

// TODO what uses this?
func (w *World) getAllPlayers() []player.Player {
	return w.playerList.All()
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

func (w *World) HandleIncomingMessage(msg *message.IncomingMessage) {
	switch messageType := msg.Request.GetMessageType(); messageType {
	case "exits":
		w.handleExits(msg)
	case "logout":
		w.handleLogout(msg)
	case "look":
		w.handleLook(msg)
	case "move":
		w.handleMove(msg)
	case "say":
		w.handleSay(msg)
	case "tell":
		w.handleTell(msg)
	case "tell_all":
		w.handleTellAll(msg)
	case "who":
		w.handleWho(msg) // TODO show connected players and the room they are in (if any)
	default:
		log.Printf("UNHANDLED messageType: %s, body %s", messageType, msg.Request)
		msg.Player.Send(message.Response{
			MessageType: messageType,
			Successful:  false,
			ResultCode:  "UNKNOWN_MESSAGE_TYPE",
		})
	}
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
