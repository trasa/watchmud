package world

import (
	"github.com/trasa/watchmud/loader"
	"log"
)

func (w *World) initialLoad(worldFilesDirectory string) (err error) {
	wb, err := loader.BuildWorld(worldFilesDirectory)
	if err != nil {
		return err
	}
	w.zones = wb.Zones
	settings := wb.Settings

	w.startRoom = w.zones[settings.StartZone].Rooms[settings.StartRoom]
	w.voidRoom = w.zones[settings.VoidZone].Rooms[settings.VoidRoom]

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
