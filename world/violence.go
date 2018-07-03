package world

import (
	"github.com/trasa/watchmud/mudtime"
	"log"
)

// Walk through all the combat going on and
// do the things to make the combat happen.
func (world *World) DoViolence(pulse mudtime.PulseCount) {

	// for each fight that is going on
	// determine if its "time" to do something
	// if so, determine what to do
	// then do it, updating the state
	// continue onwards.

	for _, fight := range world.fightLedger.GetFights() {
		// each fighter should have a speed, like fast medium slow,
		// and then we can take that into account vs the last time
		// that Violence happened - comparing it to PulseCount maybe.
		// I don't want to have the details of pulse count or real-world
		// clocks being part of mob definitions as that will make it a
		// headache to tune these settings.
		if fight.CanDoViolence(pulse) {
			log.Printf("fight now! %s", fight)
			fight.LastPulse = pulse
		}
		// TODO keep track of fights that are over
	}
	// TODO remove fights that are over from fightLedger
}
