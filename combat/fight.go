package combat

import (
	"fmt"
	"time"

	"github.com/watchmud/watchmud/rules"
)

type Fight struct {
	Fighter Combatant
	// TODO fighter speed - slow, fast, medium, whatever...
	Fightee   Combatant
	LastPulse rules.PulseCount

	// seq is the order fights started in, which the ledger -- a map -- has
	// no other way to know. See FightLedger.retarget.
	seq uint64
}

func newFight(fighter Combatant, fightee Combatant, seq uint64) *Fight {
	return &Fight{
		Fighter:   fighter,
		Fightee:   fightee,
		LastPulse: rules.PulseNever,
		seq:       seq,
	}
}

func (f *Fight) String() string {
	return fmt.Sprintf("%s fighting %s", f.Fighter.Name(), f.Fightee.Name())
}

func (f *Fight) CanDoViolence(normalViolenceDuration time.Duration, now rules.PulseCount) bool {
	// TODO mob speeds
	return rules.DurationBetween(f.LastPulse, now) >= normalViolenceDuration
}
