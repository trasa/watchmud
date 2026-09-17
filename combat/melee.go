package combat

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/trasa/watchmud/rules"
)

type MeleeAttackResult struct {
	WasHit     bool
	Damage     int
	DamageType string // TODO should be enum
	// TODO other status affects: staggering, and so on
}

func (result MeleeAttackResult) String() string {
	if result.WasHit {
		return fmt.Sprintf("Hit! %d damage!", result.Damage)
	}
	return fmt.Sprintf("Missed.")
}

// AttemptMeleeAttack calculates the result of a melee attack.
// Rolls d20, apply appropriate modifiers, if value > victim's AC, then that's a hit,
// calculate the correct damage.
func AttemptMeleeAttack(roller rules.Roller, fighter Combatant, victim Combatant) (MeleeAttackResult, error) {
	roll, err := roller.Roll("1d20")
	if err != nil {
		return MeleeAttackResult{}, err
	}
	log.Trace().Msgf("%s attacks %s, rolls a %d", fighter.Name(), victim.Name(), roll)
	// TODO need impl for critical failure and critical success
	//	criticalFailure := roll == 1
	//	criticalSuccess := roll == 20

	modifiedRoll := roll + fighter.CalculateMeleeRollModifiers()
	wasHit := modifiedRoll >= victim.ArmorClass()
	wasHitStr := "missed."
	damage := 0
	if wasHit {
		wasHitStr = "hit!"
		damage, err = calculateDamage(roller, fighter, victim)
		if err != nil {
			return MeleeAttackResult{}, err
		}
	}
	log.Trace().
		Msgf("%s attacks %s with modified roll of %d vs ac of %d, %s With %d damage.",
			fighter.Name(),
			victim.Name(),
			modifiedRoll,
			victim.ArmorClass(),
			wasHitStr,
			damage)
	return MeleeAttackResult{
		WasHit:     wasHit,
		Damage:     damage,
		DamageType: "TODO",
	}, nil
}

func calculateDamage(roller rules.Roller, fighter Combatant, victim Combatant) (int, error) {
	damage, err := roller.Roll(fighter.WeaponDamageRoll())
	if err != nil {
		return 0, err
	}
	modifier := ""

	if victim.HasResistanceTo(fighter.WeaponDamageType()) {
		modifier = " (resistance)"
		damage = damage / 2
	} else if victim.IsVulnerableTo(fighter.WeaponDamageType()) {
		modifier = " (vulnerability)"
		damage = damage * 2
	}
	log.Trace().Msgf(" %s does %d damage to %s%s", fighter.Name(), damage, victim.Name(), modifier)
	return damage, nil
}
