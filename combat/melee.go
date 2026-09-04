package combat

import (
	"fmt"

	"github.com/justinian/dice"
	"github.com/rs/zerolog/log"
)

type MeleeAttackResult struct {
	WasHit     bool
	Damage     int64
	DamageType string // TODO should be enum
	// TODO other status affects: staggering, and so on
}

func (result MeleeAttackResult) String() string {
	if result.WasHit {
		return fmt.Sprintf("Hit! %d damage!", result.Damage)
	} else {
		return fmt.Sprintf("Missed.")
	}
}

func CalculateMeleeAttack(fighter Combatant, victim Combatant) MeleeAttackResult {
	// roll d20, apply appropriate modifiers (ability modifier, proficiency bonus, others ...)
	// if value > victim's AC, that's a hit.
	// TODO: advantage, disadvantage
	roll, _, _ := dice.Roll("1d20")
	log.Debug().Msgf("%s melee attacks %s, rolls a %d", fighter.Name(), victim.Name(), roll.Int())
	return meleeAttack(fighter, victim, roll)
}

func meleeAttack(fighter Combatant, victim Combatant, roll dice.RollResult) (result MeleeAttackResult) {
	if roll.Int() == 1 {
		// critical failure!
		// TODO what to do about critical failure?
		log.Debug().Msgf("%s attacks %s but CRITICALLY FAILS!", fighter.Name(), victim.Name())
		result.WasHit = false
		return
	}
	if roll.Int() == 20 {
		// critical success!
		// TODO what to do about critical success?
		log.Debug().Msgf("%s attacks %s and CRITICALLY SUCCEEDS!", fighter.Name(), victim.Name())
		result.WasHit = true
		result.Damage = calculateDamage(fighter, victim)
		return
	}
	modifiedRoll := roll.Int() + fighter.CalculateMeleeRollModifiers()
	if modifiedRoll >= victim.ArmorClass() {
		// hit!
		log.Debug().Msgf("%s attacks %s with modified roll of %d vs ac of %d, success", fighter.Name(), victim.Name(), modifiedRoll, victim.ArmorClass())
		result.WasHit = true
		result.Damage = calculateDamage(fighter, victim)
	} else {
		// miss
		log.Debug().Msgf("%s attacks %s with modified roll of %d vs ac of %d, misses!", fighter.Name(), victim.Name(), modifiedRoll, victim.ArmorClass())
		result.WasHit = false
	}
	return result
}

func calculateDamage(fighter Combatant, victim Combatant) int64 {
	damRoll, _, _ := dice.Roll(fighter.WeaponDamageRoll())
	damage := damRoll.Int()
	log.Debug().Msgf("(raw) %s does %d damage to %s", fighter.Name(), damage, victim.Name())
	if victim.HasResistanceTo(fighter.WeaponDamageType()) {
		damage = damage / 2
		log.Debug().Msgf(" %s does %d damage to %s (resistance)", fighter.Name(), damage, victim.Name())
	} else if victim.IsVulnerableTo(fighter.WeaponDamageType()) {
		damage = damage * 2
		log.Debug().Msgf(" %s does %d damage to %s (vulnerability)", fighter.Name(), damage, victim.Name())
	}
	return int64(damage)
}
