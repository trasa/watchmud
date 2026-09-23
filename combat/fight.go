package combat

import (
	"fmt"
	"time"

	"github.com/trasa/watchmud/rules"
)

type Fight struct {
	Fighter Combatant
	// TODO fighter speed - slow, fast, medium, whatever...
	Fightee   Combatant
	LastPulse rules.PulseCount
	ZoneId    string
	RoomId    string
}

func newFight(fighter Combatant, fightee Combatant, zoneId string, roomId string) *Fight {
	return &Fight{
		fighter,
		fightee,
		rules.PulseNever,
		zoneId,
		roomId,
	}
}

func (f *Fight) String() string {
	return fmt.Sprintf("%s fighting %s", f.Fighter.Name(), f.Fightee.Name())
}

func (f *Fight) CanDoViolence(normalViolenceDuration time.Duration, now rules.PulseCount) bool {
	// TODO mob speeds
	return rules.DurationBetween(f.LastPulse, now) >= normalViolenceDuration
}
