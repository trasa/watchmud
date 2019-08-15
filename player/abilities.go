package player

import "strings"

type AbilityScore int32

type Abilities struct {
	Strength     AbilityScore `json:"str"`
	Dexterity    AbilityScore `json:"dex"`
	Constitution AbilityScore `json:"con"`
	Intelligence AbilityScore `json:"int"`
	Wisdom       AbilityScore `json:"wis"`
	Charisma     AbilityScore `json:"cha"`
}

func (a *Abilities) SetScore(attributeName string, score AbilityScore) {
	switch strings.ToLower(attributeName) {
	case "str":
		a.Strength = score
	case "dex":
		a.Dexterity = score
	case "con":
		a.Constitution = score
	case "int":
		a.Intelligence = score
	case "wis":
		a.Wisdom = score
	case "cha":
		a.Charisma = score
	}
}

// find the 'most useful' score that doesn't have a value and put this there.
// if all the scores are filled then do nothing
func (a *Abilities) FillScore(score AbilityScore) {
	if a.Constitution == AbilityScore(0) {
		a.Constitution = score
	} else if a.Dexterity == AbilityScore(0) {
		a.Dexterity = score
	} else if a.Strength == AbilityScore(0) {
		a.Strength = score
	} else if a.Intelligence == AbilityScore(0) {
		a.Intelligence = score
	} else if a.Wisdom == AbilityScore(0) {
		a.Wisdom = score
	} else if a.Charisma == AbilityScore(0) {
		a.Charisma = score
	}
}

func NewAbilities(str int32, dex int32, con int32, int int32, wis int32, cha int32) *Abilities {
	return &Abilities{
		Strength:     AbilityScore(str),
		Dexterity:    AbilityScore(dex),
		Constitution: AbilityScore(con),
		Intelligence: AbilityScore(int),
		Wisdom:       AbilityScore(wis),
		Charisma:     AbilityScore(cha),
	}
}

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
