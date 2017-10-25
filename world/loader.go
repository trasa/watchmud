package world

import (
	"github.com/trasa/watchmud/loader"
	"github.com/trasa/watchmud/mobile"
	"github.com/trasa/watchmud/object"
)

func (w *World) initialLoad() {

	w.zones = loader.BuildWorld()

	// TODO startroom and void??
	w.startRoom = w.zones["void"].Rooms["start"]
	w.voidRoom = w.zones["void"].Rooms["void"]

	// once everything is loaded, we can process the zone information
	// which says which mob instances to load and where to put them,
	// and which objects to load and where to put them
	// note that this is distinct from the building of the world
	// (reading in zone, room, object and mob definitions)
	// as this will happen throughout the runtime of the world

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
