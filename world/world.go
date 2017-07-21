package world

import (
	"github.com/trasa/watchmud/message"
	"github.com/trasa/watchmud/player"
	"log"
)

type World struct {
	Zones      map[string]*Zone
	StartRoom  *Room
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

	r := NewRoom(&sampleZone, startRoomKey, "Central Portal", "It's a boring room, with boring stuff in it.")
	sampleZone.Rooms[r.Id] = r
	w.StartRoom = r

	// north room
	northRoom := NewRoom(&sampleZone, "northRoom", "North Room", "This room is north of the start.")
	sampleZone.Rooms[northRoom.Id] = northRoom

	r.North = northRoom
	northRoom.South = r

	log.Print("World built.")
	return &w
}

func (w *World) getAllPlayers() []player.Player {
	return w.PlayerList.All()
}

func (w *World) addPlayers(player ...player.Player) {
	for _, p := range player {
		w.AddPlayer(p)
	}
}

// Add this Player to the world
// putting them in the start room
func (w *World) AddPlayer(player player.Player) {
	w.PlayerList.Add(player)
	//w.StartRoom.AddPlayer(player)
}

// TODO remove player from world

func (w *World) findPlayerByName(name string) player.Player {
	return w.PlayerList.FindByName(name)
}

func (w *World) HandleIncomingMessage(msg *message.IncomingMessage) {
	log.Printf("world incoming message: %s", msg.Request)
	switch messageType := msg.Request.GetMessageType(); messageType {
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
