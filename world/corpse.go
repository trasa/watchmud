package world

import (
	"github.com/trasa/watchmud/combat"
	"github.com/trasa/watchmud/mobile"
	"github.com/trasa/watchmud/player"
	"log"
)

func (w *World) becomeCorpse(deadCombatant combat.Combatant) {
	log.Printf("%s is dead!", deadCombatant.GetName())

	switch deadCombatant.CombatantType() {
	case combat.PlayerCombatant:
		w.becomeCorpse_Player(deadCombatant.(player.Player))
	case combat.MobileCombatant:
		w.becomeCorpse_Mobile(deadCombatant.(*mobile.Instance))
	}
}

func (w *World) becomeCorpse_Player(p player.Player) {
	// TODO
}

func (w *World) becomeCorpse_Mobile(m *mobile.Instance) {
	// TODO
}
