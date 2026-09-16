package loader

import (
	"io/fs"

	"github.com/trasa/watchmud/rules"
)

func LoadCatalog(rulesFS fs.FS) (*rules.Catalog, error) {
	species, err := readJSONFile[[]*rules.Species](rulesFS, "species.json")
	if err != nil {
		return nil, err
	}

	roles, err := readJSONFile[[]*rules.Role](rulesFS, "roles.json")
	if err != nil {
		return nil, err
	}

	armor, err := readJSONFile[rules.ArmorTypeContent](rulesFS, "armor.json")
	if err != nil {
		return nil, err
	}

	c, err := rules.NewCatalog(species, roles, armor)
	return c, err
}
