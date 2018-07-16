package world

import (
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
}

func (w *World) becomeCorpse_Mobile(m *mobile.Instance) {
	// create a corpse for the mobile instance
	// load the corpse with loot
	// remove the mobile instance
	w.removeMobile(m)
}
