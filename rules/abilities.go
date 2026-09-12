package rules

// Abilities is a character's six ability scores.
//
// Nothing contributes to them any more: a lineage is cosmetic and grants no
// bonuses, and there is no class to express a preference. Every character
// starts with the same numbers, which is the point -- see Lineage and Role.
type Abilities struct {
	Str int `json:"str"`
	Dex int `json:"dex"`
	Con int `json:"con"`
	Int int `json:"int"`
	Wis int `json:"wis"`
	Cha int `json:"cha"`
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

// StandardAbilities is the array every character starts with, highest score
// first, handed out by the priority rules above.
//
// It used to take the ability preferences of the character's class and assign
// the best numbers to them. There is no class now, so there is no preference,
// so everyone begins equal. The distribution machinery stays because letting
// a player assign their own array at creation is the obvious next step, and
// that step needs exactly this.
func StandardAbilities() Abilities {
	a := Abilities{}
	// an array of ints from highest start value to lowest (these aren't random)
	// TODO this should be a content, not a constant here...
	startScores := []int{15, 14, 13, 12, 10, 8}
	for _, score := range startScores {
		a = a.fillEmptyScoreByPriority(score)
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
