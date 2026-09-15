package rules

// Roller rolls dice and comes up with random numbers.
// Can be mocked as a 'loaded dice roller' for tests.
type Roller interface {
	Roll(notation string) (int, error)
	IntN(n int) (int, error)
}
