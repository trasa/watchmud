package world

import (
	"github.com/trasa/watchmud-message"
	"github.com/trasa/watchmud/combat"
	"github.com/trasa/watchmud/mudtime"
	"log"
)

// Walk through all the combat going on and
// do the things to make the combat happen.
//
// for each fight that is going on
// determine if its "time" to do something
// if so, determine what to do
// then do it, updating the state
// continue onwards.
func (w *World) DoViolence(pulse mudtime.PulseCount) {

	for _, fight := range w.fightLedger.GetFights() {
		if fight.Fighter.IsDead() || fight.Fightee.IsDead() {
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
			log.Printf("fight result: %s", fightResult)

			var isDead = false
			if fightResult.WasHit {
				isDead = fight.Fightee.TakeMeleeDamage(fightResult.Damage)
			}
			// tell everyone what is going on
			room, found := w.findRoomById(fight.ZoneId, fight.RoomId)
			if found {
				room.Notify(message.ViolenceNotification{
					Fighter:       fight.Fighter.GetName(),
					Fightee:       fight.Fightee.GetName(),
					SuccessfulHit: fightResult.WasHit,
					Damage:        int32(fightResult.Damage),
				})
			}

			if isDead {
				w.becomeCorpse(fight.Fightee)
				// TODO award points or other reward
				w.fightLedger.EndFight(fight.Fighter)
				// tell everybody what happened
				if found {
					room.Notify(message.DeathNotification{
						Target:   fight.Fightee.GetName(),
						IsPlayer: fight.Fightee.CombatantType() == combat.PlayerCombatant,
					})
				}
			}
		}
	}
}
