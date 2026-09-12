package rules

// A Species is the broad kind of creature: Elf, Dwarf, etc. It exists to
// group lineages for the character creation menu and for flavor, and for
// nothing else.
type Species struct {
	Id          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Lineages    []*Lineage `json:"lineages"`
}

// A Lineage is a specific kind of creature, such as High Elf or Hill Dwarf.
// Characters and mobiles reference a Lineage, never a Species directly.
//
// A lineage is cosmetic. It grants no bonuses and carries no penalties, in
// the same way that picking a gender wouldn't: every lineage is mechanically
// identical, and what a character is *good at* comes from the equipment they
// are wearing (see Role). That is deliberate -- it means a lineage can be
// changed later, by a potion or a wish or a builder command, without any of
// the rebalancing that changing a race in a traditional MUD implies.
type Lineage struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`

	// Species is wired up at load time. Excluded from JSON so marshaling
	// can't walk back up the tree.
	Species *Species `json:"-"`
}

func (l *Lineage) String() string { return l.Name }
