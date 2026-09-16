package object

import (
	"iter"
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

// All the equipped items, slot and instance, in the order a player expects to
// read them -- see rules.CompareSlots. Empty slots are skipped, so nobody has
// to remember the nil check.
func (eq *Equipment) All() iter.Seq2[rules.EquipmentSlot, *Instance] {
	return func(yield func(rules.EquipmentSlot, *Instance) bool) {
		for _, slot := range eq.slots() {
			if inst := eq.eqMap[slot]; inst != nil {
				if !yield(slot, inst) {
					return
				}
			}
		}
	}
}

// slots that currently hold something, in canonical order.
func (eq *Equipment) slots() []rules.EquipmentSlot {
	slots := slices.Collect(maps.Keys(eq.eqMap))
	slices.SortFunc(slots, rules.CompareSlots)
	return slots
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
	for slot, inst := range eq.All() {
		if inst.Matches(target) {
			return slot, inst, true
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
	//ac := 10
	//for slot, inst := range eq.eqMap {
	//	switch inst.Definition.ArmorType {
	//	}
	//}
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
	for _, inst := range eq.All() {
		for roleId, weight := range inst.Definition.RoleWeights {
			totals[roleId] += weight
		}
	}
	return totals
}
