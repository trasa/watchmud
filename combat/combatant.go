package combat

// something capable of being in combat
type Combatant interface {
	GetName() string
	// apply damage to Combatant and return if they are dead
	TakeMeleeDamage(damage int64) (isDead bool)
	// u dead?
	IsDead() bool
	// what type of combatant are you?
	CombatantType() CombatantType
	// + and - to attack roll
	CalculateMeleeRollModifiers() int
	// armor class based on intrinsic properties and equipment
	ArmorClass() int
	// resistance to damage? (resistance = dam / 2)
	HasResistanceTo(damageType DamageType) bool
	// vulnerable to damage? (vuln = dam * 2)
	IsVulnerableTo(damageType DamageType) bool
	// string of dice roll for wielded weapon (i.e. "1d6")
	WeaponDamageRoll() string
	// what type of damage does wielded weapon do
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
