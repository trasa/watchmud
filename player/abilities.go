package player

type Abilities struct {
	Strength     int
	Dexterity    int
	Constitution int
	Intelligence int
	Wisdom       int
	Charisma     int
}

func AbilityScoreModifier(a int) int {
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
