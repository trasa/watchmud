package playergenerator

import (
	"github.com/trasa/watchmud/db"
	"github.com/trasa/watchmud/player"
)

type PlayerPrototype struct {
	Race             db.RaceData
	Class            db.ClassData
	InitialAbilities *player.Abilities
}

func generatePlayerPrototype(race db.RaceData, class db.ClassData) *PlayerPrototype {

	a := generateAbilities(race, class)

	return &PlayerPrototype{
		Race:             race,
		Class:            class,
		InitialAbilities: a,
	}
}

func generateAbilities(race db.RaceData, class db.ClassData) *player.Abilities {
	a := player.Abilities{}
	startScores := []player.AbilityScore{15, 14, 13, 12, 10, 8}
	for i, score := range startScores {
		if i < len(class.AbilityPreference.Preferences) {
			a.SetScore(class.AbilityPreference.Preferences[i], score)
		} else {
			// no further class preferences
			a.FillScore(score)
		}
	}
	// apply race bonuses
	a.Strength = a.Strength + player.AbilityScore(race.StrBonus)
	a.Dexterity = a.Dexterity + player.AbilityScore(race.DexBonus)
	a.Constitution = a.Constitution + player.AbilityScore(race.ConBonus)
	a.Intelligence = a.Intelligence + player.AbilityScore(race.IntBonus)
	a.Wisdom = a.Wisdom + player.AbilityScore(race.WisBonus)
	a.Charisma = a.Charisma + player.AbilityScore(race.ChaBonus)

	return &a
}
