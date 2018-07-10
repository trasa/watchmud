package combat

type TestCombatant struct {
	Name      string
	curHealth int
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
