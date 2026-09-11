package combat

import "uuid"

// Attacker and Defender ae roles in a single attack sequence.
// They swap every round, so nothing durable about an entity
// belongs here.
type Attacker interface {
	Name() string
	CalculateMeleeRollModifiers() int
	WeaponDamageRoll() string
	WeaponDamageType() DamageType
}

type Defender interface {
	Name() string
	ArmorClass() int
	HasResistanceTo(damageType DamageType) bool
	IsVulnerableTo(damageType DamageType) bool
}

// Combatant is an entity that persists across rounds. The ledger
// and DoViolence work in these terms, melee.go does not.
type Combatant interface {
	Attacker
	Defender
	Id() uuid.UUID
	TakeMeleeDamage(damager int64) bool
	Dead() bool
}

type CombatantType int

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
