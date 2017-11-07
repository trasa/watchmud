package loader

import (
	"github.com/trasa/watchmud/direction"
	"github.com/trasa/watchmud/mobile"
	"github.com/trasa/watchmud/object"
	"github.com/trasa/watchmud/spaces"
	"github.com/trasa/watchmud/zonereset"
	"log"
	"time"
)

type WorldBuilder struct {
	zones map[string]*spaces.Zone
}

func BuildWorld() map[string]*spaces.Zone {
	worldBuilder := WorldBuilder{
		zones: make(map[string]*spaces.Zone),
	}

	worldBuilder.loadZoneManifest()
	worldBuilder.loadRooms()
	worldBuilder.loadObjectDefinitions()
	worldBuilder.loadMobileDefinitions()
	worldBuilder.loadZoneInstructions()
	return worldBuilder.zones
}

// Retrieve the zone manifest; prepare the zone objects to be
// populated by rooms, objects, mobiles (but don't process the
// zone commands yet)
func (wb *WorldBuilder) loadZoneManifest() {

	// here, we'd look up something from the database, or something.
	voidZone := spaces.NewZone("void", "The Void", zonereset.NEVER, 0)
	wb.zones[voidZone.Id] = voidZone

	zoneResetDuration := time.Minute * 3
	sampleZone := spaces.NewZone("sample", "sample zone", zonereset.ALWAYS, zoneResetDuration)
	wb.zones[sampleZone.Id] = sampleZone
}

func (wb *WorldBuilder) loadRooms() {

	// here, we'd look up something from the database...

	// The VOID. When you're not really in a room.
	voidZone := wb.zones["void"]
	voidZone.AddRoom(spaces.NewRoom(voidZone, "void", "The Void", "You see nothing but endless void."))
	// void doesn't have any directions

	// zone "sample"
	currentZone := wb.zones["sample"]
	/*
		  north -- northeast
		   |         |
		central -- east

	*/
	// central room (Start)
	currentZone.AddRoom(spaces.NewRoom(currentZone, "start", "Central Portal", "It's a boring room, with boring stuff in it."))

	// north room
	currentZone.AddRoom(spaces.NewRoom(currentZone, "northRoom", "North Room", "This room is north of the start."))

	// northeast
	currentZone.AddRoom(spaces.NewRoom(currentZone, "northeastRoom", "North East Room", "It's north, and also East."))

	// east
	currentZone.AddRoom(spaces.NewRoom(currentZone, "eastRoom", "East Room", "This room is east of the start."))

	// once all the rooms for the zones are created, we can wire the directions up
	// central <-> north
	// void.start -> void.north
	// void.north -> void.start
	wb.connectRooms("sample", "start", direction.NORTH, "sample", "northRoom")
	wb.connectRooms("sample", "northRoom", direction.SOUTH, "sample", "start")

	// central <-> east
	wb.connectRooms("sample", "start", direction.EAST, "sample", "eastRoom")
	wb.connectRooms("sample", "eastRoom", direction.WEST, "sample", "start")

	// north <-> northeast
	wb.connectRooms("sample", "northRoom", direction.EAST, "sample", "northeastRoom")
	wb.connectRooms("sample", "northeastRoom", direction.WEST, "sample", "northRoom")

	// east <-> northeast
	wb.connectRooms("sample", "eastRoom", direction.NORTH, "sample", "northeastRoom")
	wb.connectRooms("sample", "northeastRoom", direction.SOUTH, "sample", "eastRoom")

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
func (wb *WorldBuilder) loadObjectDefinitions() {

	// for each zone: create all object definitions
	// lets put "something" in the central portal room
	z := wb.zones["sample"]
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

func (wb *WorldBuilder) loadMobileDefinitions() {
	z := wb.zones["sample"]

	// walker- somebody to walk around randomly
	z.AddMobileDefinition(mobile.NewDefinition("walker",
		"walker",
		"sample",
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
		"sample",
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

func (wb *WorldBuilder) loadZoneInstructions() {
	z := wb.zones["sample"]

	// Object: put "fountain" instance in room "start", max of 1
	z.AddCommand(spaces.CreateObject{
		ObjectDefinitionId: "fountain",
		RoomId:             "start",
		InstanceMax:        1,
	})

	// Object: put "knife" instance in room "north", max of 1
	z.AddCommand(spaces.CreateObject{
		ObjectDefinitionId: "knife",
		RoomId:             "northRoom",
		InstanceMax:        1,
	})

	// Mob: put "walker" instance in room "start", max of 2
	z.AddCommand(spaces.CreateMobile{
		MobileDefinitionId: "walker",
		RoomId:             "start",
		InstanceMax:        2,
	})

	// Mob: put "scripty" instance in room "north", max of 1
	z.AddCommand(spaces.CreateMobile{
		MobileDefinitionId: "scripty",
		RoomId:             "northRoom",
		InstanceMax:        1,
	})
}
