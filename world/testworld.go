package world

import (
	"log"

	"github.com/trasa/watchmud-message/slot"
	"github.com/trasa/watchmud/combat"
	"github.com/trasa/watchmud/loader"
	"github.com/trasa/watchmud/mobile"
	"github.com/trasa/watchmud/object"
	"github.com/trasa/watchmud/player"
	"github.com/trasa/watchmud/spaces"
	"github.com/trasa/watchmud/zonereset"
)

// create a new test world
func newTestWorld() (*World, error) {

	content := loader.Content{
		Zones: map[string]*spaces.Zone{},
	}
	testZone := spaces.NewZone("sample", "Sample Zone", zonereset.NEVER, 0)
	content.Zones[testZone.Id] = testZone

	w := &World{
		content:     &content,
		playerList:  player.NewList(),
		playerRooms: NewPlayerRoomMap(),
		fightLedger: combat.NewFightLedger(),
	}
	w.initializeHandlerMap()

	testRoom := spaces.NewRoom(testZone, "start", "Test Room", "this is a test room.")
	testZone.Rooms[testRoom.Id] = testRoom
	w.StartRoom = testRoom
	w.VoidRoom = testRoom

	// put something in the start room
	knifeDefn := object.NewDefinition("knife",
		"knife",
		testZone.Id,
		object.Weapon,
		[]string{},
		"knife",
		"A knife is on the ground.",
		slot.Wield)
	knifeObj, err := object.NewInstance(knifeDefn)
	if err != nil {
		log.Printf("Error in TestWorld, failed to create knife!")
		return nil, err
	}
	testRoom.AddInventory(knifeObj)

	ironHelmetDefn := object.NewDefinition("helmet",
		"iron_helmet",
		testZone.Id,
		object.Armor,
		[]string{"helm", "iron", "helmet"},
		"iron helmet",
		"An iron helmet is on the ground.",
		slot.Head)
	ironHelmetObj, err := object.NewInstance(ironHelmetDefn)
	if err != nil {
		log.Printf("Error in TestWorld, failed to create iron helmet")
		return nil, err
	}
	testRoom.AddInventory(ironHelmetObj)

	mobDefn := mobile.NewDefinition("targetDrone", "Target Drone",
		testZone.Id,
		[]string{"target", "drone"},
		"Target Drone",
		"Target Drone buzzes around",
		25,
		mobile.WanderingDefinition{
			CanWander: false,
		},
		10)
	testZone.AddMobileDefinition(mobDefn)
	testRoom.AddMobile(mobile.NewInstance(mobDefn))

	return w, nil
}
