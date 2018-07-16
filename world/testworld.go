package world

import (
	"github.com/trasa/watchmud-message/slot"
	"github.com/trasa/watchmud/combat"
	"github.com/trasa/watchmud/mobile"
	"github.com/trasa/watchmud/object"
	"github.com/trasa/watchmud/player"
	"github.com/trasa/watchmud/spaces"
	"github.com/trasa/watchmud/zonereset"
)

// create a new test world
func newTestWorld() *World {

	testZone := spaces.NewZone("sample", "Sample Zone", zonereset.NEVER, 0)

	w := &World{
		zones:       make(map[string]*spaces.Zone),
		playerList:  player.NewList(),
		playerRooms: NewPlayerRoomMap(),
		fightLedger: combat.NewFightLedger(),
	}
	w.initializeHandlerMap()
	w.zones[testZone.Id] = testZone

	testRoom := spaces.NewRoom(testZone, "start", "Test Room", "this is a test room.")
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

	mobDefn := mobile.NewDefinition("targetDrone", "Target Drone",
		testZone.Id,
		[]string{"target", "drone"},
		"Target Drone",
		"Target Drone buzzes around",
		25,
		mobile.WanderingDefinition{
			CanWander: false,
		})
	testZone.AddMobileDefinition(mobDefn)
	testRoom.AddMobile(mobile.NewInstance(mobDefn))

	return w
}
