package combat

import "uuid"

// Attacker and Defender ae roles in a single attack sequence.
// They swap every round, so nothing durable about an entity
// belongs here.
//
// Both carry Power, since combat reads the difference between the two.
type Attacker interface {
	Name() string
	Power() int
	CalculateMeleeRollModifiers() int
	WeaponDamageRoll() string
	WeaponDamageType() DamageType
}

type Defender interface {
	Name() string
	Power() int
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
	TakeMeleeDamage(damage int) bool
	Dead() bool
}

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
