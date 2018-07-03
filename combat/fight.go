package combat

import (
	"fmt"
	"github.com/trasa/watchmud/mudtime"
)

type Fight struct {
	Fighter   Combatant
	Fightee   Combatant
	LastPulse mudtime.PulseCount
}

func NewFight(fighter Combatant, fightee Combatant) *Fight {
	return &Fight{
		fighter,
		fightee,
		mudtime.PulseCountNever,
	}
}

func (f Fight) String() string {
	return fmt.Sprintf("%s fighting %s", f.Fighter.GetName(), f.Fightee.GetName())
}

func (f *Fight) CanDoViolence(now mudtime.PulseCount) bool {
	return f.Fighter.CanDoViolence(f.LastPulse, now)
}
