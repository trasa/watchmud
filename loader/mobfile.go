package loader

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type mobFileEntry struct {
	Id                  string         `json:"id"`
	Name                string         `json:"name"`
	Aliases             []string       `json:"aliases"`
	ShortDescription    string         `json:"short_description"`
	DescriptionInRoom   string         `json:"description_in_room"`
	WanderingDefinition WanderingEntry `json:"wandering_definition"`
	Flags               []string       `json:"flags"`
	MaxHealth           int64          `json:"max_health"`
	AC int `json:"ac"`
}

type WanderingEntry struct {
	CanWander             bool     `json:"can_wander"`
	CheckFrequencySeconds int      `json:"check_frequency_seconds"`
	CheckPercentage       int      `json:"check_percentage"`
	WanderStyle           int      `json:"wander_style"`
	Path                  []string `json:"path"`
}

func readMobFile(filename string) ([]mobFileEntry, error) {
	var result []mobFileEntry
	// mob file might not exist, thats ok
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return result, nil
		} else {
			return nil, err
		}
	}
	if err = json.Unmarshal(bytes, &result); err != nil {
		return nil, err
	}
	return result, nil
}
