package dice

import (
	"math/rand/v2"

	jdice "github.com/justinian/dice"
)

type Roller struct {
	rand *rand.Rand
}

func New(seed [32]byte) *Roller {
	return &Roller{rand: rand.New(rand.NewChaCha8(seed))}
}

func (r *Roller) Roll(notation string) (int, error) {
	result, _, err := jdice.Roll(notation)
	if err != nil {
		return 0, err
	}
	return result.Int(), nil
}

func (r *Roller) IntN(n int) (int, error) {
	return r.rand.IntN(n), nil
}
