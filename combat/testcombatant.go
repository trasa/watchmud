package combat

type TestCombatant struct {
	Name      string
	curHealth int
	dead      bool
}

func NewTestCombatant(name string) *TestCombatant {
	return &TestCombatant{
		Name:      name,
		curHealth: 100, // TODO need a default
	}
}

func (t *TestCombatant) GetName() string {
	return t.Name
}

func (t *TestCombatant) TakeMeleeDamage(damage int) (isDead bool) {
	t.curHealth = t.curHealth - damage
	return t.curHealth <= 0
}

func (t *TestCombatant) IsDead() bool {
	return t.dead
}

func (t *TestCombatant) CombatantType() CombatantType {
	return NoCombatantType
}
