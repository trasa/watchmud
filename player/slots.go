package player

import (
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

func (s *Slots) IsSlotInUse(l slot.Location) (result bool) {
	_, result = s.slotMap[l]
	return
}

func (s *Slots) IsItemInUse(obj *object.Instance) bool {
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
