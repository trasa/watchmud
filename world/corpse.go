package world

import (
	"fmt"
	"github.com/trasa/watchmud-message/slot"
	"github.com/trasa/watchmud/combat"
	"github.com/trasa/watchmud/mobile"
	"github.com/trasa/watchmud/object"
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
	corpseName := fmt.Sprintf("the corpse of %s", m.Definition.Name)
	corpseDefn := object.NewDefinition("",
		corpseName,
		"",
		object.Corpse,
		m.Definition.Aliases,
		corpseName,
		fmt.Sprintf("The corpse of %s is lying here.", m.Definition.Name),
		slot.None)
	corpse := object.NewInstance(corpseDefn)

	// transfer m's possessions over to the corpse
	// TODO mobiles can't have possessions at the moment, not implemented yet..

	w.getRoomContainingMobile(m).AddInventory(corpse)

	// remove the mobile instance
	w.removeMobile(m)
}
