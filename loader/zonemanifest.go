package loader

import (
	"encoding/json"
	"io/ioutil"
)

type zoneManifestFileEntry struct {
	Id              string `json:"id"`
	Name            string `json:"name"`
	ResetMode       int    `json:"reset_mode"`
	LifetimeMinutes int    `json:"lifetime_minutes"`
}

func readZoneManifest(filename string) ([]zoneManifestFileEntry, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var result []zoneManifestFileEntry
	if err = json.Unmarshal(bytes, &result); err != nil {
		return nil, err
	}
	return result, nil
}
