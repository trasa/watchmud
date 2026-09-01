package loader

type SettingsFile struct {
	VoidZone     string `json:"void.zone"`
	VoidRoom     string `json:"void.room"`
	StartZone    string `json:"start.zone"`
	StartRoom    string `json:"start.room"`
	DonationZone string `json:"donation.zone"`
	DonationRoom string `json:"donation.room"`
}
