package world

import (
	"github.com/trasa/watchmud/object"
	"github.com/trasa/watchmud/player"
	"github.com/trasa/watchmud/spaces"
)

// create a new test world
func newTestWorld() *World {

	testZone := spaces.Zone{
		Id:    "sample",
		Name:  "Sample Zone",
		Rooms: make(map[string]*spaces.Room),
	}
	w := &World{
		zones:       make(map[string]*spaces.Zone),
		playerList:  player.NewList(),
		playerRooms: NewPlayerRoomMap(),
	}
	w.initializeHandlerMap()
	w.zones[testZone.Id] = &testZone

	testRoom := spaces.NewRoom(&testZone, "start", "Test Room", "this is a test room.")
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
	knifeObj := object.NewInstance(knifeDefn)
	testRoom.AddInventory(knifeObj)

	return w
}
