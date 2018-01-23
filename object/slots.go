package object

import (
	"github.com/trasa/watchmud/slot"
	"github.com/trasa/watchmud/thing"
)

type Slots struct {
	slotMap map[slot.Location]*Instance
	inv     Inventoryer
}

func NewSlots(inv Inventoryer) *Slots {
	return &Slots{
		slotMap: make(map[slot.Location]*Instance),
		inv:     inv,
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

func (slots *Slots) Set(s slot.Location, obj *Instance) error {
	if _, exists := slots.inv.Inventory().Get(obj.Id()); !exists {
		return &InstanceNotFoundError{Id: obj.Id()}
	}

	// is it already being used somewhere else?
	// TODO

	if err := verifyObjectForSlot(s, obj); err != nil {
		return err
	}
	slots.slotMap[s] = obj
	return nil
}

func verifyObjectForSlot(s slot.Location, obj *Instance) error {
	switch s {
	case slot.Wield:
		// is this instance valid to be a primary weapon?
		if !obj.CanEquipWeapon() {
			return &InstanceNotWeaponError{Id: obj.Id()}
		}
	}
	return nil
}

type Inventoryer interface {
	Inventory() thing.Map
}
