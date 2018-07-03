package combat

import "github.com/trasa/watchmud/mudtime"

// something capable of being in combat
type Combatant interface {
	GetName() string
	CanDoViolence(lastTime mudtime.PulseCount, now mudtime.PulseCount) bool
}
