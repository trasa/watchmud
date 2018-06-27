package object

import (
	"github.com/trasa/watchmud/slot"
)

type Slots struct {
	slotMap map[slot.Location]*Instance
}

func NewSlots() *Slots {
	return &Slots{
		slotMap: make(map[slot.Location]*Instance),
	}
}

func (slots *Slots) GetAll() map[slot.Location]*Instance {
	m := make(map[slot.Location]*Instance)
	for k, v := range slots.slotMap {
		m[k] = v
	}
	return m
}

func (slots *Slots) Get(s slot.Location) *Instance {
	if s == slot.None {
		return nil
	}
	return slots.slotMap[s]
}

func (slots *Slots) Set(s slot.Location, obj *Instance) {
	slots.slotMap[s] = obj
}

func (slots *Slots) IsSlotInUse(s slot.Location) (result bool) {
	_, result = slots.slotMap[s]
	return
}

func (slots *Slots) IsItemInUse(obj *Instance) bool {
	for _, inst := range slots.slotMap {
		if inst == obj {
			return true
		}
	}
	return false
}
