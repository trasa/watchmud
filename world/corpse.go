package world

import (
	"fmt"
	"uuid"

	"github.com/rs/zerolog/log"
	"github.com/trasa/watchmud/combat"
	"github.com/trasa/watchmud/mobile"
	"github.com/trasa/watchmud/object"
	"github.com/trasa/watchmud/rules"
)

// becomeCorpse if you are dead
func (w *World) becomeCorpse(deadCombatant combat.Combatant) {
	log.Printf("%s is dead!", deadCombatant.Name())

	// if you were fighting, you stop
	w.fightLedger.EndAllFightsWith(deadCombatant.Id())

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
		rules.ObjectCategoryCorpse,
		m.Definition.Aliases,
		corpseName,
		fmt.Sprintf("The corpse of %s is lying here.", m.Definition.Name),
		rules.SlotNone,
		// A corpse is not armor. It used to claim to be cloth, which was
		// worth nothing only because the armor table has no row for the slot
		// it can't be worn in anyway -- two accidents holding hands. Saying
		// none means it stays worth nothing if either of those changes.
		rules.ArmorTypeNone)

	corpse := object.NewInstance(uuid.New(), d)
	// TODO transfer m's possessions over to the corpse
	// TODO mobiles can't have possessions at the moment, not implemented yet..
	r := w.getRoomContainingMobile(m)
	if r == nil {
		log.Warn().Msgf("becomeMobileCorpse: could not find room containing mobile %s", m.Definition.Name)
		return
	}

	w.removeMobile(m)

	if err := r.Inventory.Add(corpse); err != nil {
		log.Error().Msgf("becomeMobileCorpse: could not add corpse %s to room %s, %v", corpse.Definition.Name, r.Name, err)
	}
}
