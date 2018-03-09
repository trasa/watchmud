package world

import (
	"github.com/trasa/watchmud/object"
	"github.com/trasa/watchmud/player"
	"github.com/trasa/watchmud/slot"
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
		object.Weapon,
		[]string{},
		"knife",
		"A knife is on the ground.",
		slot.None)
	knifeObj := object.NewInstance(knifeDefn)
	testRoom.AddInventory(knifeObj)

	ironHelmetDefn := object.NewDefinition("helmet",
		"iron_helmet",
		testZone.Id,
		object.Armor,
		[]string{"helm", "iron", "helmet"},
		"iron helmet",
		"An iron helmet is on the ground.",
		slot.Head)
	ironHelmetObj := object.NewInstance(ironHelmetDefn)
	testRoom.AddInventory(ironHelmetObj)
	return w
}
