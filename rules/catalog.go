package rules

import "fmt"

// Catalog holds the static game content loaded from the world files. It is
// built once at startup and treated as read-only thereafter.
type Catalog struct {
	Species  map[string]*Species
	Lineages map[string]*Lineage
	Classes  map[string]*Class
}

func NewCatalog(species []*Species, classes []*Class) (*Catalog, error) {
	speciesMap, lineageMap, err := indexSpecies(species)
	if err != nil {
		return nil, err
	}

	classMap, err := indexClasses(classes)
	if err != nil {
		return nil, err
	}

	return &Catalog{
		Species:  speciesMap,
		Lineages: lineageMap,
		Classes:  classMap,
	}, nil
}

// IndexSpecies wires up parent pointers and builds the lookup maps, failing
// on duplicate or missing ids. Lineage ids must be unique across all species,
// since that is what gets stored against a player.
func indexSpecies(all []*Species) (map[string]*Species, map[string]*Lineage, error) {
	species := make(map[string]*Species, len(all))
	lineages := make(map[string]*Lineage)

	for _, s := range all {
		if s.Id == "" {
			return nil, nil, fmt.Errorf("species %q: missing id", s.Name)
		}
		if prev, dup := species[s.Id]; dup {
			return nil, nil, fmt.Errorf("duplicate species id %q (%s and %s)", s.Id, prev.Name, s.Name)
		}
		species[s.Id] = s

		if len(s.Lineages) == 0 {
			return nil, nil, fmt.Errorf("species %q: has no lineages", s.Id)
		}
		for _, l := range s.Lineages {
			if l.Id == "" {
				return nil, nil, fmt.Errorf("species %q: lineage %q missing id", s.Id, l.Name)
			}
			if prev, dup := lineages[l.Id]; dup {
				return nil, nil, fmt.Errorf("duplicate lineage id %q (%s and %s)",
					l.Id, prev.Species.Id, s.Id)
			}
			l.Species = s
			lineages[l.Id] = l
		}
	}
	return species, lineages, nil
}

func indexClasses(all []*Class) (map[string]*Class, error) {
	classes := make(map[string]*Class, len(all))
	for _, c := range all {
		if c.Id == "" {
			return nil, fmt.Errorf("class %q: missing id", c.Name)
		}
		if prev, dup := classes[c.Id]; dup {
			return nil, fmt.Errorf("duplicate class id %q (%s and %s)", c.Id, prev.Name, c.Name)
		}
		classes[c.Id] = c
	}
	return classes, nil
}
