package object

import (
	"maps"
	"slices"

	"github.com/trasa/watchmud/rules"
)

type Equipment struct {
	eqMap map[rules.EquipmentSlot]*Instance
}

func NewEquipment() *Equipment {
	return &Equipment{
		eqMap: make(map[rules.EquipmentSlot]*Instance),
	}
}

func (eq *Equipment) All() map[rules.EquipmentSlot]*Instance {
	return maps.Clone(eq.eqMap)
}

func (eq *Equipment) At(slot rules.EquipmentSlot) *Instance {
	return eq.eqMap[slot]
}

func (eq *Equipment) Equip(slot rules.EquipmentSlot, obj *Instance) {
	eq.eqMap[slot] = obj
}

func (eq *Equipment) Unequip(slot rules.EquipmentSlot) {
	eq.eqMap[slot] = nil
}

func (eq *Equipment) Find(target string) (rules.EquipmentSlot, *Instance, bool) {
	for _, l := range slices.Sorted(maps.Keys(eq.eqMap)) {
		inst := eq.eqMap[l]
		if inst == nil {
			continue
		}
		if inst.Matches(target) {
			return l, inst, true
		}
	}
	return rules.SlotNone, nil, false
}

func (eq *Equipment) Equipped(slot rules.EquipmentSlot) bool {
	return eq.eqMap[slot] != nil
}

func (eq *Equipment) ItemEquipped(item *Instance) bool {
	if item == nil {
		return false
	}
	for _, inst := range eq.eqMap {
		if inst == item {
			return true
		}
	}
	return false
}

// ArmorClass is determined by summing what's equipped,
// using the ArmorType to choose the bonuses.
func (eq *Equipment) ArmorClass() int {
	// TODO
	return 10
}

// RoleWeights totals what everything equipped contributes to each role,
// keyed on rules.Role.Id. This is the whole input to a character's role:
// change what's in the slots and the totals change with it.
//
// Returned without resolving to a rules.Role, because player has no business
// holding the catalog -- the handler that needs a name calls
// rules.Catalog.RoleFor with this.
func (eq *Equipment) RoleWeights() map[string]int {
	totals := make(map[string]int)
	for _, inst := range eq.eqMap {
		if inst == nil {
			continue // a slot whose item didn't survive a content edit
		}
		// TODO work in progress
		/*
			for roleId, weight := range inst.Definition.RoleWeights {
				totals[roleId] += weight
			}*/
	}
	return totals
}
