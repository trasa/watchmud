package world

import (
	"github.com/rs/zerolog/log"
	"github.com/trasa/watchmud/combat"
	"github.com/trasa/watchmud/event"
	"github.com/trasa/watchmud/object"
	"github.com/trasa/watchmud/player"
	"github.com/trasa/watchmud/rules"
	"github.com/trasa/watchmud/spaces"
)

// wearFromBlow is the cost of a landed blow to the gear involved in it: a
// point off one piece of the defender's armor, and a point off the weapon the
// attacker swung.
//
// Only landed blows wear anything. A miss costs nobody anything, which is the
// simplest rule that still makes the fights you actually have the thing that
// wears your kit out.
//
// It lives in world/ rather than combat/ because combat works in terms of
// Attacker and Defender, which know about armor class and damage and nothing
// about equipment -- and mobs have no equipment at all. The type switch is
// here, where the concrete types already are.
func (w *World) wearFromBlow(attacker combat.Combatant, defender combat.Combatant, room *spaces.Room) {
	if p, isPlayer := defender.(*player.Player); isPlayer {
		w.wearArmor(p, room)
	}
	if p, isPlayer := attacker.(*player.Player); isPlayer {
		w.wearWeapon(p, room)
	}
}

// wearArmor takes a point off one piece of what the defender is wearing.
//
// One piece, chosen at random, rather than all of them: a character in a full
// suit would otherwise wear out five times as fast as one in a shirt, which
// is precisely backwards. There is no hit location system to ask, so the
// roller decides where the blow landed.
func (w *World) wearArmor(p *player.Player, room *spaces.Room) {
	pieces := p.Equipment().DamageableArmor()
	if len(pieces) == 0 {
		return
	}
	i, err := w.roller.IntN(len(pieces))
	if err != nil {
		log.Error().Err(err).Str("playerName", p.Name()).Msg("wearArmor: choosing where the blow landed")
		return
	}
	w.wear(p, pieces[i], 1, room)
}

// wearWeapon takes a point off whatever the attacker swung. Nothing wielded,
// or bare hands, costs nothing.
func (w *World) wearWeapon(p *player.Player, room *spaces.Room) {
	weapon := p.Equipment().At(rules.SlotWield)
	if weapon == nil || !weapon.WearsOut() || weapon.Broken() {
		return
	}
	w.wear(p, weapon, 1, room)
}

// wearFromDeath is what dying costs the gear you died in: a share of every
// piece, all at once, rather than the single point a blow costs.
//
// Everything worn, not one piece chosen at random, because there is nothing
// to choose between -- the whole kit was there. The share comes off what each
// piece started at (rules.DurabilityTable.DeathLoss), so the good armor and
// the cheap tunic cost the same number of deaths rather than the same number
// of points.
//
// This is the second caller of Instance.Damage and it needed nothing new from
// it. The pieces that give out announce themselves exactly as they do in a
// fight.
func (w *World) wearFromDeath(p *player.Player, room *spaces.Room) {
	table := w.content.Catalog.Durability

	damaged := 0
	for _, item := range p.Equipment().DamageableGear() {
		loss := table.DeathLoss(item.Definition.MaxDurability)
		if loss == 0 {
			continue
		}
		w.wear(p, item, loss, room)
		damaged++
	}

	if damaged > 0 {
		// The breaks are announced on their own. This is for the far more
		// common death, the one where everything survives: without it the
		// cost is invisible until the player thinks to check.
		p.Send(event.GearDamaged{Items: damaged})
	}
}

// wear takes the points off, and tells the room if that was the blow that
// finished the item.
func (w *World) wear(p *player.Player, item *object.Instance, points int, room *spaces.Room) {
	if !item.Damage(points) {
		return
	}
	broke := event.Broke{
		Actor: p.Name(),
		// the bare name, not the short description: this is the one message
		// that reads possessively ("your chain shirt"), and short
		// descriptions carry their own article ("a chain shirt").
		Item: item.Definition.Name,
	}
	// Breaking is worth seeing happen to somebody else, so it goes to the
	// room -- except when the fight has outlived its room, where the person
	// it happened to still needs to know.
	if room == nil {
		p.Send(broke)
		return
	}
	room.Notify(broke)
}
