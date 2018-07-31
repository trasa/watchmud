package behavior

import (
	"fmt"
	"strings"
)

type Behavior int32

//go:generate stringer -type=Behavior
const (
	None Behavior = iota
	NoTake
)

type BehaviorSet map[Behavior]bool

func NewBehaviorSet() BehaviorSet {
	return make(BehaviorSet)
}

func (bs BehaviorSet) Add(b Behavior) {
	bs[b] = true
}

func (bs BehaviorSet) Contains(b Behavior) bool {
	return bs[b]
}

func (bs BehaviorSet) ToStringList() (result []string) {
	for k, v := range bs {
		if v {
			result = append(result, k.String())
		}
	}
	return
}

func StringToBehavior(name string) (Behavior, error) {
	if name == "" {
		return None, nil
	}
	name = strings.ToUpper(name)
	stridx := strings.Index(strings.ToUpper(_Behavior_name), name)
	if stridx < 0 {
		return None, fmt.Errorf("behavior '%s' not found", name)
	}

	for pos, catidx := range _Behavior_index {
		if stridx == int(catidx) {
			return Behavior(pos), nil
		}
	}
	// shouldn't happen?
	return None, fmt.Errorf("could not find index %d for behavior '%s'", stridx, name)
}
