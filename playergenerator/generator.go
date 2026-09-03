package playergenerator

import (
	"github.com/trasa/watchmud/rules"
)

type PlayerPrototype struct {
	Lineage          rules.Lineage
	Class            rules.Class
	InitialAbilities rules.Abilities
}

func GeneratePlayerPrototype(lineage rules.Lineage, class rules.Class) *PlayerPrototype {
	a := assignStandardArray(class.AbilityPreference)

	return &PlayerPrototype{
		Lineage:          lineage,
		Class:            class,
		InitialAbilities: a,
	}
}

// assignStandardArray distributes the standard array over the six abilities in
// the given preference order; anything unnamed falls through to FillScore
func assignStandardArray(order []string) rules.Abilities {
	a := rules.Abilities{}
	// an array of ints from highest start value to lowest (these aren't random)
	// TODO this should be a rule, not a constant here...
	startScores := []int{15, 14, 13, 12, 10, 8}
	// staring with the highest value, map the value to the ability listed
	// first (second, third...) in the class.AbilityPreference. once we're
	// past that number of scores, the rest just get set in order.
	for i, score := range startScores {
		if i < len(order) {
			a = a.Set(order[i], score)
		} else {
			// no further class preferences
			a = a.FillEmptyScoreByPriority(score)
		}
	}

	return a
}
