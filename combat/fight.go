package combat

import "fmt"

type Fight struct {
	Fighter Combatant
	Fightee Combatant
	// TODO: timing for last attack
}

func NewFight(fighter Combatant, fightee Combatant) *Fight {
	return &Fight {
		fighter,
		fightee,
	}
}

func (f Fight) String() string {
	return fmt.Sprintf("%s fighting %s", f.Fighter.GetName(), f.Fightee.GetName())
}