package world

import (
	"log"
	"os"

	"github.com/trasa/watchmud/loader"
	"github.com/trasa/watchmud/serverconfig"
)

func (w *World) initialLoad(cfg serverconfig.Config) (err error) {
	wb := loader.NewWorldBuilder()
	if err := wb.Load(os.DirFS(cfg.WorldFilesDir)); err != nil {
		return err
	}
	w.zones = wb.Zones
	settings := wb.Settings

	w.StartRoom = w.zones[settings.StartZone].Rooms[settings.StartRoom]
	w.VoidRoom = w.zones[settings.VoidZone].Rooms[settings.VoidRoom]

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
