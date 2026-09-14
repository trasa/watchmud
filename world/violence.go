package world

import (
	"github.com/rs/zerolog/log"
	"github.com/trasa/watchmud/combat"
	"github.com/trasa/watchmud/event"
	"github.com/trasa/watchmud/mudtime"
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
			fightResult := combat.CalculateMeleeAttack(fight.Fighter, fight.Fightee)
			log.Debug().Msgf("fight result: %s", fightResult)

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
				w.becomeCorpse(fight.Fightee)
				// TODO award points or other reward
				w.fightLedger.EndFight(fight.Fighter)
				// tell everybody what happened
				if found {
					room.Notify(event.Died{
						Target: fight.Fightee.Name(),
					})
				}
			}
		}
	}
}
