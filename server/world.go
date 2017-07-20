package server

import (
	"github.com/trasa/watchmud/message"
	"github.com/trasa/watchmud/player"
	"github.com/trasa/watchmud/response"
	"log"
)

type World struct {
	Zones              map[string]*Zone
	StartRoom          *Room
	knownPlayersByName map[string]player.Player // TODO move to players?
	PlayerList         *player.List
}

var startZoneKey = "sample"
var startRoomKey = "start"

// Constructor for World
func NewWorld() *World {
	// Build a very boring world.
	w := World{
		Zones:              make(map[string]*Zone),
		knownPlayersByName: make(map[string]player.Player),
		PlayerList:         player.NewList(),
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

// Add this Player to the world
// putting them in the start room
func (w *World) addPlayer(player player.Player) {
	w.knownPlayersByName[player.GetName()] = player // TODO players
	w.PlayerList.Add(player)
	//w.StartRoom.AddPlayer(player)
}

// TODO remove player

func (w *World) findPlayerByName(name string) player.Player {
	return w.knownPlayersByName[name]
}

func (w *World) handleIncomingMessage(message *message.IncomingMessage) {
	log.Printf("world incoming message: %s", message.Body)
	switch messageType := message.Body["msg_type"]; messageType {
	case "login":
		log.Printf("login received: %s", message.Body)
		w.handleLogin(message)
	case "tell":
		log.Printf("tell: %s", message.Body)
		w.handleTell(message)
	case "tell_all":
		log.Printf("Tell All: %s", message)
		w.handleTellAll(message)
	default:
		log.Printf("UNHANDLED messageType: %s, body %s", messageType, message.Body)
		message.Client.Send(response.Response{
			MessageType: messageType,
			Successful:  false,
			ResultCode:  "UNKNOWN_MESSAGE_TYPE",
		})
	}
}

func (w *World) SendToAllPlayers(message interface{}) {
	w.PlayerList.Iter(func(p player.Player) {
		p.Send(message)
	})
}
