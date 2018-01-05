package loader

import (
	"encoding/json"
	"io/ioutil"
)

type SettingsFile struct {
	VoidZone     string `json:"void.zone"`
	VoidRoom     string `json:"void.room"`
	StartZone    string `json:"start.zone"`
	StartRoom    string `json:"start.room"`
	DonationZone string `json:"donation.zone"`
	DonationRoom string `json:"donation.room"`
}

func readSettingsFile(filename string) (*SettingsFile, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	result := &SettingsFile{}
	if err = json.Unmarshal(bytes, result); err != nil {
		return nil, err
	}
	return result, nil
}
