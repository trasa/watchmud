package world

import (
	"log"
	"time"

	"github.com/trasa/watchmud/zonereset"
)

func (w *World) DoZoneActivity() {
	w.doZoneActivity(time.Now())
}

func (w *World) doZoneActivity(now time.Time) {
	// for each zone,
	// do whatever needs to happen on pulse
	// for example, a zone reset
	for _, z := range w.content.Zones {
		if z.ResetMode == zonereset.NO_PLAYERS || z.ResetMode == zonereset.ALWAYS {
			// is it time yet for this zone's lifetime?
			if now.Sub(z.LastReset) > z.Lifetime {
				if errs := z.Reset(w.mobileRooms); len(errs) != 0 {
					log.Printf("World.DoZoneActivity Errors: :%s", errs)
				}
			}
		}
	}
}
