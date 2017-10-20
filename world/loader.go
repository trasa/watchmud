package world

import (
	"github.com/trasa/watchmud/mobile"
	"github.com/trasa/watchmud/object"
	"time"
)

func (w *World) initialLoad() {

	sampleZone := Zone{
		Id:    "sample",
		Name:  "Sample Zone",
		Rooms: make(map[string]*Room),
	}
	w.zones[sampleZone.Id] = &sampleZone

	/*
		  north -- northeast
		   |         |
		central -- east

	*/
	// central room (Start)
	centralPortalRoom := NewRoom(&sampleZone, "start", "Central Portal", "It's a boring room, with boring stuff in it.")
	sampleZone.Rooms[centralPortalRoom.Id] = centralPortalRoom
	w.startRoom = centralPortalRoom

	// north room
	northRoom := NewRoom(&sampleZone, "northRoom", "North Room", "This room is north of the start.")
	sampleZone.Rooms[northRoom.Id] = northRoom

	// northeast
	northeastRoom := NewRoom(&sampleZone, "northeastRoom", "North East Room", "It's north, and also East.")
	sampleZone.Rooms[northeastRoom.Id] = northeastRoom

	// east
	eastRoom := NewRoom(&sampleZone, "eastRoom", "East Room", "This room is east of the start.")
	sampleZone.Rooms[eastRoom.Id] = eastRoom

	// central <-> north
	centralPortalRoom.North = northRoom
	northRoom.South = centralPortalRoom

	// central <-> east
	centralPortalRoom.East = eastRoom
	eastRoom.West = centralPortalRoom

	// north <-> northeast
	northRoom.East = northeastRoom
	northeastRoom.West = northRoom

	// east <-> northeast
	eastRoom.North = northeastRoom
	northeastRoom.South = eastRoom

	// The VOID. When you're not really in a room.
	w.voidRoom = NewRoom(nil, "void", "The Void", "You see nothing but endless void.")

	// lets put "something" in the central portal room
	fountainDefn := object.NewDefinition(
		"fountain",
		"fountain",
		object.OTHER,
		[]string{"fount"},
		"fountain",
		"A fountain bubbles quietly.")

	// TODO instance ids
	fountainObj := object.NewInstance("fountain", fountainDefn)
	// put the obj in the room
	centralPortalRoom.AddInventory(fountainObj)

	// that's not a knife....wait, yes it is.
	knifeDefn := object.NewDefinition(
		"knife",
		"knife",
		object.WEAPON,
		[]string{},
		"knife",
		"A knife is on the ground.")
	// TODO instance ids
	knifeObj := object.NewInstance("knife", knifeDefn)
	// knife is in room
	centralPortalRoom.AddInventory(knifeObj)

	// walker- somebody to walk around randomly
	walkerDefn := mobile.NewDefinition("walker", "walker", []string{}, "The Walker walks.",
		"The walker stands here...for now.",
		mobile.WanderingDefinition{
			CanWander:       true,
			CheckFrequency:  time.Second * 30, // check for movement every N seconds
			CheckPercentage: 0.50,             // % chance of movement on each check
			Style:           mobile.WANDER_RANDOM,
		})

	// TODO instance ids
	walkerObj := mobile.NewInstance("walker", walkerDefn)
	w.AddMobiles(centralPortalRoom, walkerObj)

	// scripty -- scripted action in a mob
	scriptyDefn := mobile.NewDefinition("scripty", "scripty", []string{}, "Scripty thinks about things.",
		"Scripty is pondering something.",
		mobile.WanderingDefinition{
			CanWander:       true,
			CheckFrequency:  time.Second * 10,
			CheckPercentage: 1.0,
			Style:           mobile.WANDER_FOLLOW_PATH,
			Path:            []string{"start", "northRoom", "northeastRoom", "eastRoom"},
		})
	//TODO instance ids
	scriptyObj := mobile.NewInstance("scripty", scriptyDefn)
	w.AddMobiles(centralPortalRoom, scriptyObj)
}
