package loader

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type instructionFileEntry struct {
	Type        string `json:"type"`
	ObjectId    string `json:"object_id"`
	MobileId    string `json:"mobile_id"`
	RoomId      string `json:"room_id"`
	InstanceMax int    `json:"instance_max"`
}

func readInstructionFile(filename string) ([]instructionFileEntry, error) {
	var result []instructionFileEntry
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
