package world

import (
	"fmt"
	"uuid"

	"github.com/rs/zerolog/log"
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

	switch c := deadCombatant.(type) {
	case *mobile.Instance:
		w.becomeMobileCorpse(c)
		//case *player.Player:
		//	TODO!
		// default - nothing happens
	}
}

// BecomeMobileCorpse turns a mobile instance into a corpse.
// The corpse has the loot from the mobile instance.
func (w *World) becomeMobileCorpse(m *mobile.Instance) {
	// create a corpse for the mobile instance
	// load the corpse with loot
	corpseName := fmt.Sprintf("the corpse of %s", m.Definition.Name)
	d := object.NewDefinition("",
		corpseName,
		"",
		object.Corpse,
		m.Definition.Aliases,
		corpseName,
		fmt.Sprintf("The corpse of %s is lying here.", m.Definition.Name),
		slot.None)

	corpse := object.NewInstance(uuid.New(), d)
	// TODO transfer m's possessions over to the corpse
	// TODO mobiles can't have possessions at the moment, not implemented yet..
	w.removeMobile(m)
	// figure out what room to put the corpse into
	r := w.getRoomContainingMobile(m)
	if r == nil {
		log.Warn().Msgf("becomeMobileCorpse: could not find room containing mobile %s", m.Definition.Name)
		return
	}
	if err := r.Inventory.Add(corpse); err != nil {
		log.Error().Msgf("becomeMobileCorpse: could not add corpse %s to room %s, %v", corpse.Definition.Name, r.Name, err)
	}
}
