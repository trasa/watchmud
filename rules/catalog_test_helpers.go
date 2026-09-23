package rules

import "time"

func newMudTime() MudTime {
	return MudTime{
		Mobile:   time.Second * 10,
		Violence: time.Second,
		Zone:     time.Minute,
	}
}

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
		{Id: "tank", Name: "Tank", FromArmor: true},
		{Id: "healer", Name: "Healer"},
		{Id: "striker", Name: "Striker"},
	}
}

// NewTestArmor is deliberately empty: most tests are about what gear declares
// by hand, and a table here would quietly add weight to every one of them.
// Tests about armor build their own.
func NewTestArmor() ArmorTypeContent {
	return make(ArmorTypeContent)
}

func NewTestCatalog() (*Catalog, error) {
	return NewCatalog(
		newMudTime(),
		newTestSpecies(),
		NewTestRoles(),
		NewTestArmor(),
	)
}
