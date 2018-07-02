package world

import "log"

// Walk through all the combat going on and
// do the things to make the combat happen.
func (world *World) DoViolence() {

	// for each fight that is going on
	// determine if its "time" to do something
	// if so, determine what to do
	// then do it, updating the state
	// continue onwards.

	for _, fight := range world.fightLedger.GetFights() {
		log.Printf("%s", fight)
	}


}
