package combat

import "github.com/trasa/watchmud/mudtime"

type TestCombatant struct {
	Name string
}

func NewTestCombatant(name string) *TestCombatant {
	return &TestCombatant{
		Name: name,
	}
}

func (t *TestCombatant) GetName() string {
	return t.Name
}

func (t *TestCombatant) CanDoViolence(last mudtime.PulseCount, now mudtime.PulseCount) bool {
	// TODO
	return false
}
