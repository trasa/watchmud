package combat

import (
	"fmt"
	"github.com/trasa/watchmud/mudtime"
	"time"
)

type Fight struct {
	Fighter Combatant
	// TODO fighter speed - slow, fast, medium, whatever...
	Fightee   Combatant
	LastPulse mudtime.PulseCount
	ZoneId    string
	RoomId    string
}

func NewFight(fighter Combatant, fightee Combatant, zoneId string, roomId string) *Fight {
	return &Fight{
		fighter,
		fightee,
		mudtime.PulseCountNever,
		zoneId,
		roomId,
	}
}

func (f Fight) String() string {
	return fmt.Sprintf("%s fighting %s", f.Fighter.GetName(), f.Fightee.GetName())
}

func (f *Fight) CanDoViolence(now mudtime.PulseCount) bool {
	// for now. lets say that clientplayer fight speed is "normal"
	// and that means they fight every 3 seconds.
	return mudtime.DurationBetween(f.LastPulse, now) >= time.Duration(time.Second*3)
}
