package world

import (
	"github.com/trasa/watchmud/message"
	"github.com/trasa/watchmud/player"
	"log"
)

//noinspection GoNameStartsWithPackageName
type World struct {
	Zones      map[string]*Zone
	StartRoom  *Room
	VoidRoom   *Room
	PlayerList *player.List
}

var startZoneKey = "sample"
var startRoomKey = "start"

// Constructor for World
func New() *World {
	// Build a very boring world.
	w := World{
		Zones:      make(map[string]*Zone),
		PlayerList: player.NewList(),
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

// TODO do we need this? merge into AddPlayer()?
func (w *World) addPlayers(player ...player.Player) {
	for _, p := range player {
		w.AddPlayer(p)
	}
}

// Add this Player to the world
// putting them in the start room
func (w *World) AddPlayer(p player.Player) {
	w.PlayerList.Add(p)
	w.StartRoom.AddPlayer(p)
}

func (w *World) EnterRoom(p player.Player, r *Room) {
	// TODO
	// tell everybody about the new player in the room
}

func (w *World) LeaveRoom(p player.Player, r *Room) {
	// TODO
	// tell everybody about the player who left the room
}

func (w *World) GetRoomContainingPlayer(p player.Player) *Room {
	// TODO have to figure out this structure
	if w.StartRoom.PlayerList.FindByName(p.GetName()) != nil {
		return w.StartRoom
	} else {
		return nil
	}
}

// TODO remove player from world

func (w *World) findPlayerByName(name string) player.Player {
	return w.PlayerList.FindByName(name)
}

func (w *World) HandleIncomingMessage(msg *message.IncomingMessage) {
	switch messageType := msg.Request.GetMessageType(); messageType {
	case "go":
		w.handleGo(msg)
	case "look":
		w.handleLook(msg)
	case "tell":
		w.handleTell(msg)
	case "tell_all":
		w.handleTellAll(msg)
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
