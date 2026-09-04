package combat

// Combatant is something capable of being in combat
type Combatant interface {
	Name() string
	TakeMeleeDamage(damage int64) (isDead bool)
	Dead() bool
	Type() CombatantType
	CalculateMeleeRollModifiers() int
	ArmorClass() int
	HasResistanceTo(damageType DamageType) bool
	IsVulnerableTo(damageType DamageType) bool
	WeaponDamageRoll() string
	WeaponDamageType() DamageType
}

type CombatantType int

const (
	NoCombatantType CombatantType = iota // for testing
	PlayerCombatant
	MobileCombatant
)

type DamageType int

const (
	Acid DamageType = iota
	Bludgeoning
	Cold
	Fire
	Force
	Lightning
	Necrotic
	Piercing
	Poison
	Psychic
	Radiant
	Slashing
	Thunder
)
