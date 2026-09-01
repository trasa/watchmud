package loader

import "io/fs"

type Settings struct {
	VoidZone     string `json:"void.zone"`
	VoidRoom     string `json:"void.room"`
	StartZone    string `json:"start.zone"`
	StartRoom    string `json:"start.room"`
	DonationZone string `json:"donation.zone"`
	DonationRoom string `json:"donation.room"`
}

func LoadSettings(worldFS fs.FS) (*Settings, error) {
	settings := Settings{}
	if err := loadInto(&settings, worldFS, "settings.json"); err != nil {
		return nil, err
	}
	return &settings, nil
}
