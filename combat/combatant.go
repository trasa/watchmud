package combat

// something capable of being in combat
type Combatant interface {
	GetName() string
	TakeMeleeDamage(damage int) (isDead bool)
}
