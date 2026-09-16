package rules

func newTestSpecies() []*Species {
	l := Lineage{
		Id:   "human",
		Name: "human",
	}
	s := Species{
		Id:   "human",
		Name: "human",
	}
	l.Species = &s
	s.Lineages = []*Lineage{&l}
	return []*Species{&s}
}

// NewTestRoles mirrors the shape of content/rules/roles.json: three roles, in
// a fixed declaration order, since that order is what breaks ties.
func NewTestRoles() []*Role {
	return []*Role{
		{Id: "tank", Name: "Tank"},
		{Id: "healer", Name: "Healer"},
		{Id: "striker", Name: "Striker"},
	}
}

func NewTestArmor() ArmorTypeContent {
	return make(ArmorTypeContent)
}

func NewTestCatalog() (*Catalog, error) {
	return NewCatalog(newTestSpecies(), NewTestRoles(), NewTestArmor())
}
