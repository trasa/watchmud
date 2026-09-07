package rules

func NewTestSpecies() []*Species {
	humanBonus := Abilities{1, 1, 1, 1, 1, 1}
	l := Lineage{
		Id:         "human",
		Name:       "human",
		OwnBonuses: Abilities{},
	}
	s := Species{
		Id:         "human",
		Name:       "human",
		OwnBonuses: humanBonus,
	}
	l.Species = &s
	s.Lineages = []*Lineage{&l}
	return []*Species{&s}
}

func NewTestClasses() []*Class {
	c := Class{
		Id:                "fighter",
		Name:              "fighter",
		AbilityPreference: []string{"str", "dex", "con"},
	}
	return []*Class{&c}
}

func NewTestCatalog() (*Catalog, error) {
	species := NewTestSpecies()
	classes := NewTestClasses()
	c, err := NewCatalog(species, classes)
	return c, err
}
