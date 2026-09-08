package world

import (
	"time"

	"github.com/rs/zerolog/log"
	"github.com/trasa/watchmud/zonereset"
)

// DoZoneActivity for each zone based on time.Now()
// For example, zone resets.
func (w *World) DoZoneActivity() {
	w.doZoneActivity(time.Now())
}

// doZoneActivity for each zone based on pulse time
// For example, zone resets.
func (w *World) doZoneActivity(now time.Time) {
	for _, z := range w.content.Zones {
		if z.ResetMode == zonereset.NO_PLAYERS || z.ResetMode == zonereset.ALWAYS {
			// is it time yet for this zone's lifetime?
			if now.Sub(z.LastReset) > z.Lifetime {
				if errs := z.Reset(w.mobileRooms); len(errs) != 0 {
					for _, err := range errs {
						log.Warn().Str("zone", z.Id).Err(err).Msg("zone reset error")
					}
				}
			}
		}
	}
}
