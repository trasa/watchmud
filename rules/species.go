package rules

// A Species is the broad kind of creature: Elf, Dwarf, etc. It carries the
// traits shared by every member, however specialized their lineage.
type Species struct {
	Id         string     `json:"id""`
	Name       string     `json:"name"`
	OwnBonuses Abilities  `json:"bonuses"`
	Lineages   []*Lineage `json:"lineages"`
}

// A Lineage is a specific species of creature, such as High Elf or Hill Dwarf.
// Characters and mobiles reference a Lineage, never a Species directly.
type Lineage struct {
	Id         string    `json:"id"`
	Name       string    `json:"name"`
	OwnBonuses Abilities `json:"bonuses"`

	// Species is wired up at load time. Excluded from JSON so marshaling
	// can't walk back up the tree.
	Species *Species `json:"-"`
}

// Bonus returns everything the lineage contributes to a character's
// abilities: what's true of the species, plus what's additionally true of
// this lineage. Nothing overrides; contributions accumulate.
func (l *Lineage) Bonus() Abilities {
	if l.Species == nil {
		return l.OwnBonuses
	}
	return l.Species.OwnBonuses.Add(l.OwnBonuses)
}

func (l *Lineage) String() string { return l.Name }
