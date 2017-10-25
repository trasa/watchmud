package loader

import (
	"github.com/trasa/watchmud/direction"
	"github.com/trasa/watchmud/mobile"
	"github.com/trasa/watchmud/object"
	"github.com/trasa/watchmud/spaces"
	"log"
	"time"
)

type WorldBuilder struct {
	zones map[string]*spaces.Zone
}

func BuildWorld() map[string]*spaces.Zone {
	worldBuilder := WorldBuilder{
		make(map[string]*spaces.Zone),
	}

	worldBuilder.loadZoneManifest()
	worldBuilder.loadRooms()
	worldBuilder.loadObjectDefinitions()
	worldBuilder.loadMobileDefinitions()
	return worldBuilder.zones
}

// Retrieve the zone manifest; prepare the zone objects to be
// populated by rooms, objects, mobiles (but don't process the
// zone commands yet)
func (worldBuilder *WorldBuilder) loadZoneManifest() {

	// here, we'd look up something from the database, or something.
	sampleZone := spaces.NewZone("void", "void zone")
	worldBuilder.zones[sampleZone.Id] = sampleZone
}

func (wb *WorldBuilder) loadRooms() {

	// here, we'd look up something from the database...

	// TODO get rid of void room?
	// The VOID. When you're not really in a room.
	//w.voidRoom = spaces.NewRoom(nil, "void", "The Void", "You see nothing but endless void.")

	// zone "void"
	currentZone := wb.zones["void"]
	/*
		  north -- northeast
		   |         |
		central -- east

	*/
	// central room (Start)
	centralPortalRoom := spaces.NewRoom(currentZone, "start", "Central Portal", "It's a boring room, with boring stuff in it.")
	currentZone.AddRoom(centralPortalRoom)
	// TODO some better way of indicating the start room from configuration
	//w.startRoom = centralPortalRoom

	// north room
	northRoom := spaces.NewRoom(currentZone, "northRoom", "North Room", "This room is north of the start.")
	currentZone.AddRoom(northRoom)

	// northeast
	northeastRoom := spaces.NewRoom(currentZone, "northeastRoom", "North East Room", "It's north, and also East.")
	currentZone.AddRoom(northeastRoom)

	// east
	eastRoom := spaces.NewRoom(currentZone, "eastRoom", "East Room", "This room is east of the start.")
	currentZone.AddRoom(eastRoom)

	// once all the rooms for the zones are created, we can wire the directions up
	// central <-> north
	// void.start -> void.north
	// void.north -> void.start
	wb.connectRooms("void", "start", direction.NORTH, "void", "northRoom")
	wb.connectRooms("void", "northRoom", direction.SOUTH, "void", "start")

	// central <-> east
	wb.connectRooms("void", "start", direction.EAST, "void", "eastRoom")
	wb.connectRooms("void", "eastRoom", direction.WEST, "void", "start")

	// north <-> northeast
	wb.connectRooms("void", "northRoom", direction.EAST, "void", "northeastRoom")
	wb.connectRooms("void", "northeastRoom", direction.WEST, "void", "northRoom")

	// east <-> northeast
	wb.connectRooms("void", "eastRoom", direction.NORTH, "void", "northeastRoom")
	wb.connectRooms("void", "northeastRoom", direction.SOUTH, "void", "eastRoom")

}

func (wb *WorldBuilder) connectRooms(sourceZoneId string, sourceRoomId string, dir direction.Direction, destZoneId string, destRoomId string) {
	sourceZone := wb.zones[sourceZoneId]
	if sourceZone == nil {
		log.Printf("ConnectRooms failed: sourceZoneId '%s' not found", sourceZoneId)
		return
	}
	destZone := wb.zones[destZoneId]
	if destZone == nil {
		log.Printf("ConnectRooms failed: destZoneId '%s' not found", destZoneId)
		return
	}
	sourceRoom := sourceZone.Rooms[sourceRoomId]
	log.Printf("source room is %s", sourceRoom)
	if sourceRoom == nil {
		log.Printf("ConnectRooms failed: zone %s sourceRoomId '%s' not found", sourceZoneId, sourceRoomId)
		return
	}
	destRoom := destZone.Rooms[destRoomId]
	if destRoom == nil {
		log.Printf("ConnectRooms failed: zone %s destRoomId '%s' not found", destZoneId, destRoomId)
		return
	}
	sourceRoom.Set(dir, destRoom)
}
func (w *WorldBuilder) loadObjectDefinitions() {

	// for each zone: create all object definitions
	// lets put "something" in the central portal room
	z := w.zones["void"]
	z.AddObjectDefinition(object.NewDefinition(
		"fountain",
		"fountain",
		z.Id,
		object.OTHER,
		[]string{"fount"},
		"fountain",
		"A fountain bubbles quietly."))

	// that's not a knife....wait, yes it is.
	z.AddObjectDefinition(object.NewDefinition(
		"knife",
		"knife",
		z.Id,
		object.WEAPON,
		[]string{},
		"knife",
		"A knife is on the ground."))

}

func (w *WorldBuilder) loadMobileDefinitions() {
	z := w.zones["void"]

	// walker- somebody to walk around randomly
	z.AddMobileDefinition(mobile.NewDefinition("walker",
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
	z.AddMobileDefinition(mobile.NewDefinition("scripty",
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
