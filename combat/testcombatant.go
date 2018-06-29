package combat

type TestCombatant struct {
	x int // don't be empty or instance comparison doesn't work
}

func NewTestCombatant() *TestCombatant {
	return &TestCombatant{}
}
