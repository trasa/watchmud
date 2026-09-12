package player

import (
	"maps"
	"slices"

	"github.com/trasa/watchmud/object"
	"github.com/trasa/watchmud/slot"
)

type Slots struct {
	slotMap map[slot.Location]*object.Instance
}

func NewSlots() *Slots {
	return &Slots{
		slotMap: make(map[slot.Location]*object.Instance),
	}
}

func (s *Slots) GetAll() map[slot.Location]*object.Instance {
	m := make(map[slot.Location]*object.Instance)
	for k, v := range s.slotMap {
		m[k] = v
	}
	return m
}

func (s *Slots) Get(l slot.Location) *object.Instance {
	if l == slot.None {
		return nil
	}
	return s.slotMap[l]
}

func (s *Slots) Set(l slot.Location, obj *object.Instance) {
	s.slotMap[l] = obj
}

// Clear empties a slot. It deletes the key rather than storing a nil, because
// a nil left behind reads as "occupied" to anything that tests for presence
// -- which is how a removed item could make its slot unusable forever.
func (s *Slots) Clear(l slot.Location) {
	delete(s.slotMap, l)
}

// FindByNameOrAlias looks for an equipped item by what the player called it,
// searching slots in order so the same input always finds the same item.
func (s *Slots) FindByNameOrAlias(target string) (slot.Location, *object.Instance, bool) {
	for _, l := range slices.Sorted(maps.Keys(s.slotMap)) {
		inst := s.slotMap[l]
		if inst == nil {
			continue
		}
		if inst.Definition.Name == target || inst.Definition.HasAlias(target) {
			return l, inst, true
		}
	}
	return slot.None, nil, false
}

// IsSlotInUse tests for something actually in the slot. A nil value is not
// "in use": see Clear.
func (s *Slots) IsSlotInUse(l slot.Location) bool {
	return s.slotMap[l] != nil
}

func (s *Slots) IsItemInUse(obj *object.Instance) bool {
	if obj == nil {
		return false
	}
	for _, inst := range s.slotMap {
		if inst == obj {
			return true
		}
	}
	return false
}

func (s *Slots) ArmorClass() int {
	// TODO build me
	return 10
}

// RoleWeights totals what everything equipped contributes to each role,
// keyed on rules.Role.Id. This is the whole input to a character's role:
// change what's in the slots and the totals change with it.
//
// Returned without resolving to a rules.Role, because player has no business
// holding the catalog -- the handler that needs a name calls
// rules.Catalog.RoleFor with this.
func (s *Slots) RoleWeights() map[string]int {
	totals := make(map[string]int)
	for _, inst := range s.slotMap {
		if inst == nil {
			continue // a slot whose item didn't survive a content edit
		}
		for roleId, weight := range inst.Definition.RoleWeights {
			totals[roleId] += weight
		}
	}
	return totals
}
