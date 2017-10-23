package world

import (
	"github.com/trasa/watchmud/object"
	"github.com/trasa/watchmud/player"
)

// create a new test world
func newTestWorld() *World {

	testZone := Zone{
		Id:    "sample",
		Name:  "Sample Zone",
		Rooms: make(map[string]*Room),
	}
	w := &World{
		zones:       make(map[string]*Zone),
		playerList:  player.NewList(),
		playerRooms: NewPlayerRoomMap(),
	}
	w.initializeHandlerMap()
	w.zones[testZone.Id] = &testZone

	testRoom := NewRoom(&testZone, "start", "Test Room", "this is a test room.")
	testZone.Rooms[testRoom.Id] = testRoom
	w.startRoom = testRoom

	// put something in the start room
	knifeDefn := object.NewDefinition("knife",
		"knife",
		testZone.Id,
		object.WEAPON,
		[]string{},
		"knife",
		"A knife is on the ground.")
	knifeObj := object.NewInstance("knife", knifeDefn)
	testRoom.AddInventory(knifeObj)

	return w
}
