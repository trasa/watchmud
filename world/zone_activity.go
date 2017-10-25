package world

import "log"

func (w *World) DoZoneActivity() {
	// for each zone,
	// do whatever needs to happen on pulse
	// for example, a zone reset
	for _, z := range w.zones {
		if err := z.DoZoneActivity(); err != nil {
			log.Printf("World.DoZoneActivity Error: %s", err)
		}
	}
}
