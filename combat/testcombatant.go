package combat

type TestCombatant struct {
	Name          string
	curHealth     int64
	dead          bool
	ac            int
	resistance    map[DamageType]bool
	vulnerability map[DamageType]bool
}

func NewTestCombatant(name string, ac int, resists []DamageType, vulnerableTo []DamageType) *TestCombatant {
	tc := &TestCombatant{
		Name:          name,
		curHealth:     100, // TODO need a default
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

func (t *TestCombatant) GetName() string {
	return t.Name
}

func (t *TestCombatant) TakeMeleeDamage(damage int64) (isDead bool) {
	t.curHealth = t.curHealth - damage
	return t.curHealth <= 0
}

func (t *TestCombatant) IsDead() bool {
	return t.dead
}

func (t *TestCombatant) CombatantType() CombatantType {
	return NoCombatantType
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
