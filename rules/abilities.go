package rules

import "strings"

// Abilities is a set of six ability scores. It serves double duty: as a
// character's actual scores, and as the bonuses a species or lineage
// contributes to them.
type Abilities struct {
	Str int `json:"str"`
	Dex int `json:"dex"`
	Con int `json:"con"`
	Int int `json:"int"`
	Wis int `json:"wis"`
	Cha int `json:"cha"`
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

func (a Abilities) Set(name string, score int) Abilities {
	switch strings.ToLower(name) {
	case "str":
		a.Str = score
	case "dex":
		a.Dex = score
	case "con":
		a.Con = score
	case "int":
		a.Int = score
	case "wis":
		a.Wis = score
	case "cha":
		a.Cha = score
	}
	return a
}

// fillEmptyScoreByPriority fills the first empty ability score with the given score.
// Where the priority is determined by rules.
// TODO: in the future we want to encode this in rules/*.json not hardcoded here.
func (a Abilities) fillEmptyScoreByPriority(score int) Abilities {
	if a.Con == 0 {
		a.Con = score
	} else if a.Dex == 0 {
		a.Dex = score
	} else if a.Str == 0 {
		a.Str = score
	} else if a.Int == 0 {
		a.Int = score
	} else if a.Wis == 0 {
		a.Wis = score
	} else if a.Cha == 0 {
		a.Cha = score
	}
	return a
}

// StandardAbilities are the standard set of ability numbers,
// distributed by assigning the preferences given first, then the rest
// by the priority rules.
func StandardAbilities(preferences []string) Abilities {
	a := Abilities{}
	// an array of ints from highest start value to lowest (these aren't random)
	// TODO this should be a content, not a constant here...
	startScores := []int{15, 14, 13, 12, 10, 8}
	// staring with the highest value, map the value to the ability listed
	// first (second, third...) in the class.AbilityPreference. once we're
	// past that number of scores, the rest just get set in order.
	for i, score := range startScores {
		if i < len(preferences) {
			a = a.Set(preferences[i], score)
		} else {
			// no further class preferences
			a = a.fillEmptyScoreByPriority(score)
		}
	}

	return a
}

/*
TODO move to rules
func AbilityScoreModifier(a AbilityScore) int {
	switch {
	case a <= 1:
		return -5
	case a <= 3:
		return -4
	case a <= 5:
		return -3
	case a <= 7:
		return -2
	case a <= 9:
		return -1
	case a <= 11:
		return 0
	case a <= 13:
		return 1
	case a <= 15:
		return 2
	case a <= 17:
		return 3
	case a <= 19:
		return 4
	case a <= 21:
		return 5
	case a <= 23:
		return 6
	case a <= 25:
		return 7
	case a <= 27:
		return 8
	case a <= 29:
		return 9
	default:
		return 10
	}
}
*/
