package world

import (
	"github.com/trasa/watchmud/direction"
	"github.com/trasa/watchmud/mobile"
	"github.com/trasa/watchmud/object"
	"time"
)

func (w *World) initialLoad() {

	// master list of zones to load (depending on server mode)
	w.loadZoneManifest()

	w.loadRooms()

	w.loadObjectDefinitions()

	w.loadMobileDefinitions()

	// once everything is loaded, we can process the zone information
	// which says which mob instances to load and where to put them,
	// and which objects to load and where to put them

	// TODO instance ids
	fountainObj := object.NewInstance("fountain", w.zones["void"].ObjectDefinitions["fountain"])
	// put the obj in the room
	w.zones["void"].Rooms["start"].AddInventory(fountainObj)

	// TODO instance ids
	knifeObj := object.NewInstance("knife", w.zones["void"].ObjectDefinitions["knife"])
	// knife is in room
	w.zones["void"].Rooms["start"].AddInventory(knifeObj)

	// TODO instance ids
	walkerObj := mobile.NewInstance("walker", w.zones["void"].MobileDefinitions["walker"])
	w.AddMobiles(w.zones["void"].Rooms["start"], walkerObj)

	//TODO instance ids
	scriptyObj := mobile.NewInstance("scripty", w.zones["void"].MobileDefinitions["scripty"])
	w.AddMobiles(w.zones["void"].Rooms["start"], scriptyObj)
}

// Retrieve the zone manifest; prepare the zone objects to be
// populated by rooms, objects, mobiles (but don't process the
// zone commands yet)
func (w *World) loadZoneManifest() {

	// here, we'd look up something from the database, or something.
	sampleZone := NewZone("void", "void zone")
	w.zones[sampleZone.Id] = sampleZone
}

// Retrieve the room information for the world, creating the
// room pointers and putting them in the world.
func (w *World) loadRooms() {

	// here, we'd look up something from the database...

	// TODO get rid of void room?
	// The VOID. When you're not really in a room.
	w.voidRoom = NewRoom(nil, "void", "The Void", "You see nothing but endless void.")

	// zone "void"
	currentZone := w.zones["void"]
	/*
		  north -- northeast
		   |         |
		central -- east

	*/
	// central room (Start)
	centralPortalRoom := NewRoom(currentZone, "start", "Central Portal", "It's a boring room, with boring stuff in it.")
	currentZone.addRoom(centralPortalRoom)
	// TODO some better way of indicating the start room from configuration
	w.startRoom = centralPortalRoom

	// north room
	northRoom := NewRoom(currentZone, "northRoom", "North Room", "This room is north of the start.")
	currentZone.addRoom(northRoom)

	// northeast
	northeastRoom := NewRoom(currentZone, "northeastRoom", "North East Room", "It's north, and also East.")
	currentZone.addRoom(northeastRoom)

	// east
	eastRoom := NewRoom(currentZone, "eastRoom", "East Room", "This room is east of the start.")
	currentZone.addRoom(eastRoom)

	// once all the rooms for the zones are created, we can wire the directions up
	// central <-> north
	// void.start -> void.north
	// void.north -> void.start
	w.ConnectRooms("void", "start", direction.NORTH, "void", "northRoom")
	w.ConnectRooms("void", "northRoom", direction.SOUTH, "void", "start")

	// central <-> east
	w.ConnectRooms("void", "start", direction.EAST, "void", "eastRoom")
	w.ConnectRooms("void", "eastRoom", direction.WEST, "void", "start")

	// north <-> northeast
	w.ConnectRooms("void", "northRoom", direction.EAST, "void", "northeastRoom")
	w.ConnectRooms("void", "northeastRoom", direction.WEST, "void", "northRoom")

	// east <-> northeast
	w.ConnectRooms("void", "eastRoom", direction.NORTH, "void", "northeastRoom")
	w.ConnectRooms("void", "northeastRoom", direction.SOUTH, "void", "eastRoom")

}

func (w *World) loadObjectDefinitions() {

	// for each zone: create all object definitions
	// lets put "something" in the central portal room
	z := w.zones["void"]
	z.addObjectDefinition(object.NewDefinition(
		"fountain",
		"fountain",
		z.Id,
		object.OTHER,
		[]string{"fount"},
		"fountain",
		"A fountain bubbles quietly."))

	// that's not a knife....wait, yes it is.
	z.addObjectDefinition(object.NewDefinition(
		"knife",
		"knife",
		z.Id,
		object.WEAPON,
		[]string{},
		"knife",
		"A knife is on the ground."))

}

func (w *World) loadMobileDefinitions() {
	z := w.zones["void"]

	// walker- somebody to walk around randomly
	z.addMobileDefinition(mobile.NewDefinition("walker",
		"walker",
		"void",
		[]string{},
		"The Walker walks.",
		"The walker stands here...for now.",
		mobile.WanderingDefinition{
			CanWander:       true,
			CheckFrequency:  time.Second * 30, // check for movement every N seconds
			CheckPercentage: 0.50,             // % chance of movement on each check
			Style:           mobile.WANDER_RANDOM,
		}))

	// scripty -- scripted action in a mob
	z.addMobileDefinition(mobile.NewDefinition("scripty",
		"scripty",
		"void",
		[]string{},
		"Scripty thinks about things.",
		"Scripty is pondering something.",
		mobile.WanderingDefinition{
			CanWander:       true,
			CheckFrequency:  time.Second * 10,
			CheckPercentage: 1.0,
			Style:           mobile.WANDER_FOLLOW_PATH,
			Path:            []string{"start", "northRoom", "northeastRoom", "eastRoom"},
		}))
}
