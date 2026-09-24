package loader

import "github.com/trasa/watchmud/rules"

type zoneManifestEntry struct {
	Id              string `json:"id"`
	Name            string `json:"name"`
	ResetMode       int    `json:"reset_mode"`
	LifetimeMinutes int    `json:"lifetime_minutes"`
	Enabled         bool   `json:"enabled"`
	// Power is the band the zone is built for; absent is the bottom band.
	Power *rules.PowerBand `json:"power"`
}
