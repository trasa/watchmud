package equip

import (
	"github.com/trasa/watchmud/object"
	"github.com/trasa/watchmud/thing"
)

type Slot int32

//go:generate stringer -type=Slot
const (
	NONE Slot = iota
	PRIMARY_HAND
	SECONDARY_HAND
)

type Slots struct {
	slotMap map[Slot]*object.Instance
	inv     Inventoryer
}

func NewSlots() *Slots {
	return &Slots{
		slotMap: make(map[Slot]*object.Instance),
	}
}

func (slots Slots) Get(s Slot) *object.Instance {
	if s == NONE {
		return nil
	}
	return slots.slotMap[s]
}

func (slots Slots) Set(s Slot, obj *object.Instance) error {
	if _, exists := slots.inv.Inventory().Get(obj.Id()); !exists {
		return &object.InstanceNotFoundError{Id: obj.Id()}
	}

	// is it already being used somewhere else?
	// TODO

	err := verifyObjectForSlot(s, obj)
	if err != nil {
		return err
	}
	slots.slotMap[s] = obj
	return nil
}

func verifyObjectForSlot(s Slot, obj *object.Instance) error {
	switch s {
	case PRIMARY_HAND, SECONDARY_HAND:
		// is this instance valid to be a primary weapon?
		if !obj.CanEquipWeapon() {
			return &object.InstanceNotWeaponError{Id: obj.Id()}
		}
	}
	return nil
}

type Inventoryer interface {
	Inventory() thing.Map
}
