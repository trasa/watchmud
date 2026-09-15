package world

import "fmt"

type loadedDice struct {
	vals []int
	i    int
}

func newLoadedDice() *loadedDice {
	return &loadedDice{}
}
func (l *loadedDice) Add(i int) {
	l.vals = append(l.vals, i)
}

func (l *loadedDice) Load(vals []int) {
	l.vals = vals
	l.i = 0
}

func (l *loadedDice) Reset() {
	l.i = 0
}

func (l *loadedDice) Roll(notation string) (int, error) {
	if l.i >= len(l.vals) {
		return 0, fmt.Errorf("loaded dice ran out of values")
	}
	l.i++
	return l.vals[l.i-1], nil
}

func (l *loadedDice) IntN(n int) (int, error) {
	if l.i >= len(l.vals) {
		return 0, fmt.Errorf("loaded dice ran out of values")
	}
	l.i++
	return l.vals[l.i-1], nil
}
