package server

import "log"

type World struct {
	Zones            map[string]*Zone
	StartRoom        *Room
	knownPlayersById map[string]*Player
}

var startZoneKey = "sample"
var startRoomKey = "start"

// Constructor for World
func NewWorld() *World {
	// Build a very boring world.
	w := World{
		Zones:            make(map[string]*Zone),
		knownPlayersById: make(map[string]*Player),
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

func (w *World) getAllPlayers() (result []*Player) {
	for _, v := range w.knownPlayersById {
		result = append(result, v)
	}
	return result
}

// Add this Player to the world
// putting them in the start room
func (w *World) AddPlayer(player *Player) {
	w.knownPlayersById[player.Id] = player
	w.StartRoom.AddPlayer(player)
}
