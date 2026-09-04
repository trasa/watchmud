package world

import (
	"fmt"
	"log"

	"github.com/trasa/watchmud-message/slot"
	"github.com/trasa/watchmud/combat"
	"github.com/trasa/watchmud/mobile"
	"github.com/trasa/watchmud/object"
)

// becomeCorpse if you are dead
func (w *World) becomeCorpse(deadCombatant combat.Combatant) {
	// TODO: should be *Combatant?
	log.Printf("%s is dead!", deadCombatant.Name())

	// if you were fighting, you stop
	w.fightLedger.EndFight(deadCombatant)

	// TODO error handling, names
	// TODO case where its a PlayerCombatant
	switch deadCombatant.Type() {
	case combat.MobileCombatant:
		w.becomeCorpse_Mobile(deadCombatant.(*mobile.Instance))
	}
}

func (w *World) becomeCorpse_Mobile(m *mobile.Instance) error {
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

	if corpse, err := object.NewInstance(corpseDefn); err != nil {
		return err
	} else {
		// transfer m's possessions over to the corpse
		// TODO mobiles can't have possessions at the moment, not implemented yet..

		w.getRoomContainingMobile(m).AddInventory(corpse)

		// remove the mobile instance
		w.removeMobile(m)
		return nil
	}
}
