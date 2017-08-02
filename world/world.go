package world

import (
	"github.com/trasa/watchmud/direction"
	"github.com/trasa/watchmud/message"
	"github.com/trasa/watchmud/player"
	"log"
)

//noinspection GoNameStartsWithPackageName
type World struct {
	Zones       map[string]*Zone
	StartRoom   *Room
	VoidRoom    *Room
	PlayerList  *player.List
	PlayerRooms *PlayerRoomMap
}

var startZoneKey = "sample"
var startRoomKey = "start"

// Constructor for World
func New() *World {
	// Build a very boring world.
	w := World{
		Zones:       make(map[string]*Zone),
		PlayerList:  player.NewList(),
		PlayerRooms: NewPlayerRoomMap(),
	}
	sampleZone := Zone{
		Id:    startZoneKey,
		Name:  "Sample Zone",
		Rooms: make(map[string]*Room),
	}
	w.Zones[sampleZone.Id] = &sampleZone

	centralPortalRoom := NewRoom(&sampleZone, startRoomKey, "Central Portal", "It's a boring room, with boring stuff in it.")
	sampleZone.Rooms[centralPortalRoom.Id] = centralPortalRoom
	w.StartRoom = centralPortalRoom

	// north room
	northRoom := NewRoom(&sampleZone, "northRoom", "North Room", "This room is north of the start.")
	sampleZone.Rooms[northRoom.Id] = northRoom

	// north room and central portal connect to each other
	centralPortalRoom.North = northRoom
	northRoom.South = centralPortalRoom

	// The VOID. When you're not really in a room.
	w.VoidRoom = NewRoom(nil, "void", "The Void", "You see nothing but endless void.")

	log.Print("World built.")
	return &w
}

// TODO what uses this?
func (w *World) getAllPlayers() []player.Player {
	return w.PlayerList.All()
}

// Add Player(s) to the world putting them in the start room,
// Don't send room notifications.
func (w *World) AddPlayer(players ...player.Player) {
	for _, p := range players {
		log.Printf("Adding Player: %s", p.GetName())
		w.PlayerList.Add(p)
		w.PlayerRooms.Add(p, w.StartRoom)
		w.StartRoom.Add(p)
	}
}

func (w *World) RemovePlayer(players ...player.Player) {
	for _, p := range players {
		log.Printf("Removing Player: %s", p.GetName())
		w.PlayerList.Remove(p)
		w.PlayerRooms.Remove(p)
	}
}

// Player is moving from src room to dest room.
func (w *World) movePlayer(p player.Player, dir direction.Direction, src *Room, dest *Room) {
	src.Leave(p, dir)
	dest.Enter(p)
	w.PlayerRooms.Add(p, dest)
}

func (w *World) getRoomContainingPlayer(p player.Player) *Room {
	return w.PlayerRooms.Get(p)
}

func (w *World) findPlayerByName(name string) player.Player {
	return w.PlayerList.FindByName(name)
}

func (w *World) HandleIncomingMessage(msg *message.IncomingMessage) {
	switch messageType := msg.Request.GetMessageType(); messageType {
	case "logout":
		w.handleLogout(msg)
	case "look":
		w.handleLook(msg)
	case "move":
		w.handleMove(msg)
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
	w.PlayerList.Iter(func(p player.Player) {
		p.Send(message)
	})
}

// Send a message to all players in the world except for one.
func (w *World) SendToAllPlayersExcept(exception player.Player, message interface{}) {
	w.PlayerList.Iter(func(p player.Player) {
		if exception != p {
			p.Send(message)
		}
	})
}
