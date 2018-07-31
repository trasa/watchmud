package combat

// something capable of being in combat
type Combatant interface {
	GetName() string
	TakeMeleeDamage(damage int64) (isDead bool)
	IsDead() bool
	CombatantType() CombatantType
}

type CombatantType int

const (
	NoCombatantType CombatantType = iota // for testing
	PlayerCombatant
	MobileCombatant
)
