package rules

import "fmt"

// Catalog holds the static game content loaded from the world files. It is
// built once at startup and treated as read-only thereafter.
type Catalog struct {
	Species  map[string]*Species
	Lineages map[string]*Lineage
	Roles    map[string]*Role

	// declaration order, preserved from the content files: the creation menu
	// reads speciesOrder, and roleOrder breaks ties in RoleFor. Iterating the
	// maps instead would shuffle both from run to run.
	speciesOrder []*Species
	roleOrder    []*Role
}

func NewCatalog(species []*Species, roles []*Role) (*Catalog, error) {
	speciesMap, lineageMap, err := indexSpecies(species)
	if err != nil {
		return nil, err
	}

	roleMap, roleOrder, err := indexRoles(roles)
	if err != nil {
		return nil, err
	}

	return &Catalog{
		Species:      speciesMap,
		Lineages:     lineageMap,
		Roles:        roleMap,
		speciesOrder: species,
		roleOrder:    roleOrder,
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

// SpeciesList returns every species in the order the content declared them,
// each still holding its own lineages in their declared order. This is the
// character creation menu.
func (c *Catalog) SpeciesList() []*Species {
	return c.speciesOrder
}

// RoleList returns every role in the order the content declared them.
func (c *Catalog) RoleList() []*Role {
	return c.roleOrder
}

// DefaultLineage is what a character gets when nobody picked one, or when the
// one they picked is no longer in the content. Since lineage is cosmetic, an
// unrecognized one is never worth refusing a login over.
func (c *Catalog) DefaultLineage() *Lineage {
	for _, s := range c.speciesOrder {
		if len(s.Lineages) > 0 {
			return s.Lineages[0]
		}
	}
	return nil
}

// RoleFor picks the role a set of equipment adds up to: the highest total
// wins, and a tie goes to whichever role the content declared first. Gear
// that contributes nothing to any role -- or no gear at all -- is no role,
// reported as nil, because "you are wearing nothing in particular" is a real
// answer and inventing a default would hide it.
func (c *Catalog) RoleFor(weights map[string]int) *Role {
	var best *Role
	bestTotal := 0
	for _, r := range c.roleOrder {
		if total := weights[r.Id]; total > bestTotal {
			best, bestTotal = r, total
		}
	}
	return best
}
