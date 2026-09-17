package world

import (
	"github.com/rs/zerolog/log"
	"github.com/trasa/watchmud/combat"
	"github.com/trasa/watchmud/event"
	"github.com/trasa/watchmud/mudtime"
	"github.com/trasa/watchmud/spaces"
)

// DoViolence walks through all the combat going on and
// makes the combat happen. For each fight, determine if
// it is "time" to do something, and if so determine what to do.
// Update the state, and continue.
func (w *World) DoViolence(pulse mudtime.PulseCount) {

	for _, fight := range w.fightLedger.GetFights() {
		if fight.Fighter.Dead() || fight.Fightee.Dead() {
			continue
		}

		// each fighter should have a speed, like fast medium slow,
		// and then we can take that into account vs the last time
		// that Violence happened - comparing it to PulseCount.
		// I don't want to have the details of pulse count or real-world
		// clocks being part of mob definitions as that will make it a
		// headache to tune these settings.
		if fight.CanDoViolence(pulse) {
			fight.LastPulse = pulse
			fightResult, err := combat.AttemptMeleeAttack(w.roller, fight.Fighter, fight.Fightee)
			if err != nil {
				log.Error().Err(err).Msg("failed to attempt melee attack")
				continue
			}
			var isDead = false
			if fightResult.WasHit {
				isDead = fight.Fightee.TakeMeleeDamage(fightResult.Damage)
			}
			// tell everyone what is going on
			room, found := w.findRoomById(fight.ZoneId, fight.RoomId)
			if found {
				room.Notify(event.Struck{
					Attacker: fight.Fighter.Name(),
					Target:   fight.Fightee.Name(),
					Hit:      fightResult.WasHit,
					Damage:   int(fightResult.Damage),
				})
			}

			if isDead {
				w.combatantDied(fight.Fightee, room, found)
				// TODO award points or other reward
			}
		}
	}
}

// combatantDied cleans up after the one who died and tells the room.
//
// Only the dead one leaves the ledger. This used to end every fight the
// *winner* was in as well, which meant killing one of two attackers quietly
// took you out of the fight with the other one -- combat with more than one
// attacker could never happen. becomeCorpse already ends the dead one's
// fights, in both directions, which is the whole of what should end here.
func (w *World) combatantDied(dead combat.Combatant, room *spaces.Room, roomFound bool) {
	w.becomeCorpse(dead)
	if roomFound {
		room.Notify(event.Died{
			Target: dead.Name(),
		})
	}
}
