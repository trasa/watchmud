package loader

import (
	"encoding/json"
	"io/ioutil"
)

type roomFileEntry struct {
	Id          string                `json:"id"`
	Name        string                `json:"name"`
	Description string                `json:"description"`
	Flags       []string              `json:"flags"`
	Exits       []*roomFileEntryExits `json:"exits"`
}

type roomFileEntryExits struct {
	Direction         int    `json:"direction"`
	DestinationZoneId string `json:"dest_zone"`
	DestinationRoomId string `json:"dest_room"`
}

func readRoomFile(filename string) ([]roomFileEntry, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var result []roomFileEntry
	if err = json.Unmarshal(bytes, &result); err != nil {
		return nil, err
	}
	return result, nil
}
