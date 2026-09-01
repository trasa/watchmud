package rules

// Abilities is a set of six ability scores. It serves double duty: as a
// character's actual scores, and as the bonuses a species or lineage
// contributes to them.
type Abilities struct {
	Str int32 `json:"str"`
	Dex int32 `json:"dex"`
	Con int32 `json:"con"`
	Int int32 `json:"int"`
	Wis int32 `json:"wis"`
	Cha int32 `json:"cha"`
}

// Add returns the component-wise sum of two ability sets.
func (a Abilities) Add(b Abilities) Abilities {
	return Abilities{
		Str: a.Str + b.Str,
		Dex: a.Dex + b.Dex,
		Con: a.Con + b.Con,
		Int: a.Int + b.Int,
		Wis: a.Wis + b.Wis,
		Cha: a.Cha + b.Cha,
	}
}
