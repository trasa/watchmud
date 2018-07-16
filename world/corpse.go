package world

import (
	"github.com/trasa/watchmud-message"
	"github.com/trasa/watchmud/combat"
	"github.com/trasa/watchmud/mobile"
	"github.com/trasa/watchmud/player"
	"log"
)

func (w *World) becomeCorpse(deadCombatant combat.Combatant) {
	log.Printf("%s is dead!", deadCombatant.GetName())

	// if you were fighting, you stop
	w.fightLedger.EndFight(deadCombatant)

	switch deadCombatant.CombatantType() {
	case combat.PlayerCombatant:
		w.becomeCorpse_Player(deadCombatant.(player.Player))
	case combat.MobileCombatant:
		w.becomeCorpse_Mobile(deadCombatant.(*mobile.Instance))
	}
}

func (w *World) becomeCorpse_Player(p player.Player) {
	// TODO
	// create a corpse? nah. too old fashioned. keep your stuff.

	// but you do get bumped back to login

	// should we force disconnect? or have some sort of landing area?
	// that is still TBD

	// tell others in room that player is dead
	playerRoom := w.getRoomContainingPlayer(p)
	if playerRoom != nil {
		playerRoom.Notify(message.DeathNotification{
			Target:   p.GetName(),
			IsPlayer: true,
		})
	}
}

func (w *World) becomeCorpse_Mobile(m *mobile.Instance) {
	// tell mobile that they are dead
	// tell room that mobile is dead
	mobileRoom := w.getRoomContainingMobile(m)
	if mobileRoom != nil {
		mobileRoom.Notify(message.DeathNotification{
			Target:   m.GetName(),
			IsPlayer: false,
		})
	}

	// create a corpse for the mobile instance
	// load the corpse with loot

	// remove the mobile instance
	w.removeMobile(m)
}
