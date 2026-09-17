package testdice

import "fmt"

type LoadedDice struct {
	vals []int
	i    int
}

func New() *LoadedDice {
	return &LoadedDice{}
}

func (l *LoadedDice) Add(i int) {
	l.vals = append(l.vals, i)
}

func (l *LoadedDice) Load(vals []int) {
	l.vals = vals
	l.i = 0
}

func (l *LoadedDice) Reset() {
	l.i = 0
}

func (l *LoadedDice) Roll(notation string) (int, error) {
	if l.i >= len(l.vals) {
		return 0, fmt.Errorf("loaded dice ran out of values")
	}
	l.i++
	return l.vals[l.i-1], nil
}

func (l *LoadedDice) IntN(n int) (int, error) {
	if l.i >= len(l.vals) {
		return 0, fmt.Errorf("loaded dice ran out of values")
	}
	l.i++
	return l.vals[l.i-1], nil
}
