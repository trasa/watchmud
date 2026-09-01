package loader

type zoneManifestEntry struct {
	Id              string `json:"id"`
	Name            string `json:"name"`
	ResetMode       int    `json:"reset_mode"`
	LifetimeMinutes int    `json:"lifetime_minutes"`
	Enabled         bool   `json:"enabled"`
}
