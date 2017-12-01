package world

import (
	"github.com/trasa/watchmud/loader"
	"log"
)

func (w *World) initialLoad(worldFilesDirectory string) (err error) {
	w.zones, err = loader.BuildWorld(worldFilesDirectory)
	if err != nil {
		return err
	}
	settings := loader.LoadWorldSettings()

	w.startRoom = w.zones[settings["start.zone"]].Rooms[settings["start.room"]]
	w.voidRoom = w.zones[settings["void.zone"]].Rooms[settings["void.room"]]

	// once everything is loaded, we can process the zone information
	// which says which mob instances to load and where to put them,
	// and which objects to load and where to put them
	// note that this is distinct from the building of the world
	// (reading in zone, room, object and mob definitions)
	// as this will happen throughout the runtime of the world

	for zoneId, z := range w.zones {
		if errs := z.Reset(w.mobileRooms); len(errs) != 0 {
			log.Printf("Error running initial zone reset for %s: %s", zoneId, errs)
		}
	}
	return nil
}
