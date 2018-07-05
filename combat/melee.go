package combat

import (
	"fmt"
	"math/rand"
	"time"
)

type MeleeAttackResult struct {
	WasHit bool
	Damage int
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
	// always 75% for now.
	return meleeAttack(fighter, victim, rand.New(rand.NewSource(time.Now().UnixNano())), 0.75)
}

func meleeAttack(fighter Combatant, victim Combatant, r *rand.Rand, percentToHit float64) (result MeleeAttackResult) {
	// figure out % for fighter to hit
	// adjust by defense abilities of victim
	// did they hit?
	chance := r.Float64()
	if chance < percentToHit {
		// yes - determine damage
		result.WasHit = true
		// for now, just assume 5 points.
		result.Damage = 5
	} else {
		// no
		result.WasHit = false
	}

	return result
}
