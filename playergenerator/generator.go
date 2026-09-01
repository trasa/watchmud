package playergenerator

import (
	"github.com/trasa/watchmud/player"
	"github.com/trasa/watchmud/rules"
)

type PlayerPrototype struct {
	Lineage          rules.Lineage
	Class            rules.Class
	InitialAbilities *player.Abilities // TODO revisit this
}

func GeneratePlayerPrototype(lineage rules.Lineage, class rules.Class) *PlayerPrototype {
	a := generateAbilities(lineage, class)

	return &PlayerPrototype{
		Lineage:          lineage,
		Class:            class,
		InitialAbilities: a,
	}
}

func generateAbilities(lineage rules.Lineage, class rules.Class) *player.Abilities {
	a := player.Abilities{}
	startScores := []player.AbilityScore{15, 14, 13, 12, 10, 8}
	for i, score := range startScores {
		if i < len(class.AbilityPreference) {
			a.SetScore(class.AbilityPreference[i], score)
		} else {
			// no further class preferences
			a.FillScore(score)
		}
	}
	// apply race bonuses
	// TODO revisit this, we probably don't want to overwrite the base values with the modifiers?
	/*a.Strength = a.Strength + player.AbilityScore(race.StrBonus)
	a.Dexterity = a.Dexterity + player.AbilityScore(race.DexBonus)
	a.Constitution = a.Constitution + player.AbilityScore(race.ConBonus)
	a.Intelligence = a.Intelligence + player.AbilityScore(race.IntBonus)
	a.Wisdom = a.Wisdom + player.AbilityScore(race.WisBonus)
	a.Charisma = a.Charisma + player.AbilityScore(race.ChaBonus)
	*/
	return &a
}
