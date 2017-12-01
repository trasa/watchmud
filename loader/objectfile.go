package loader

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type objectFileEntry struct {
	Id                  string   `json:"id"`
	Name                string   `json:"name"`
	Category            int      `json:"category"`
	Aliases             []string `json:"aliases"`
	ShortDescription    string   `json:"short_description"`
	DescriptionOnGround string   `json:"description_on_ground"`
}

func readObjectFile(filename string) ([]objectFileEntry, error) {
	var result []objectFileEntry
	// object file might not exist, and thats OK
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
