package combat

import "uuid"

type TestCombatant struct {
	id            uuid.UUID
	name          string
	curHealth     int64
	dead          bool
	ac            int
	resistance    map[DamageType]bool
	vulnerability map[DamageType]bool
}

func NewTestCombatant(name string, ac int, resists []DamageType, vulnerableTo []DamageType) *TestCombatant {
	tc := &TestCombatant{
		id:            uuid.New(),
		name:          name,
		curHealth:     100,
		ac:            ac,
		resistance:    make(map[DamageType]bool),
		vulnerability: make(map[DamageType]bool),
	}
	for _, dt := range resists {
		tc.resistance[dt] = true
	}
	for _, dt := range vulnerableTo {
		tc.vulnerability[dt] = true
	}
	return tc
}

func (t *TestCombatant) Id() uuid.UUID {
	return t.id
}

func (t *TestCombatant) Name() string {
	return t.name
}

func (t *TestCombatant) TakeMeleeDamage(damage int64) (isDead bool) {
	t.curHealth = t.curHealth - damage
	return t.curHealth <= 0
}

func (t *TestCombatant) Dead() bool {
	return t.dead
}

func (t *TestCombatant) CalculateMeleeRollModifiers() int {
	// no modifiers yet
	return 0
}

func (t *TestCombatant) ArmorClass() int {
	return t.ac
}

func (t *TestCombatant) HasResistanceTo(damageType DamageType) bool {
	_, contained := t.resistance[damageType]
	return contained
}

func (t *TestCombatant) IsVulnerableTo(damageType DamageType) bool {
	_, contained := t.vulnerability[damageType]
	return contained
}

func (t *TestCombatant) WeaponDamageRoll() string {
	return "1d6"
}

func (t *TestCombatant) WeaponDamageType() DamageType {
	return Piercing
}
