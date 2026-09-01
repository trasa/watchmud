package loader

import (
	"io/fs"

	"github.com/trasa/watchmud/rules"
)

func LoadRulesCatalog(rulesFS fs.FS) (*rules.Catalog, error) {
	species, err := readJSONFile[[]*rules.Species](rulesFS, "species.json")
	if err != nil {
		return nil, err
	}

	classes, err := readJSONFile[[]*rules.Class](rulesFS, "classes.json")
	if err != nil {
		return nil, err
	}

	c, err := rules.NewCatalog(species, classes)
	return c, err
}
