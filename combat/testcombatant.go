package combat

type TestCombatant struct {
	Name string
}

func (t *TestCombatant) GetName() string {
	return t.Name
}

func NewTestCombatant(name string) *TestCombatant {
	return &TestCombatant{
		Name: name,
	}
}
