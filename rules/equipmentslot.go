package rules

import (
	"cmp"
	"errors"
	"fmt"
	"slices"
	"strings"
)

type EquipmentSlot string

const (
	SlotNone      EquipmentSlot = ""
	SlotWield     EquipmentSlot = "wield" // item can be wielded (i.e. weapon)
	SlotHold      EquipmentSlot = "hold"  // item can be held ('hold' command)
	SlotHead      EquipmentSlot = "head"
	SlotNeck      EquipmentSlot = "neck"
	SlotBody      EquipmentSlot = "body"
	SlotAboutBody EquipmentSlot = "about_body" // like a cloak
	SlotLegs      EquipmentSlot = "legs"
	SlotFeet      EquipmentSlot = "feet"
	SlotArms      EquipmentSlot = "arms"
	SlotWrist     EquipmentSlot = "wrist"
	SlotHands     EquipmentSlot = "hands"
	SlotFingers   EquipmentSlot = "fingers"
	SlotWaist     EquipmentSlot = "waist"
)

func (e *EquipmentSlot) UnmarshalText(b []byte) error {
	slot, err := ParseEquipmentSlot(string(b))
	if err != nil {
		return err
	}
	*e = slot
	return nil
}

var ErrUnknownEquipmentSlot = errors.New("unknown equipment slot")

var equipmentSlots = []EquipmentSlot{
	SlotNone, // some equipment cannot be worn at all, so this is a valid parse target
	SlotWield,
	SlotHold,
	SlotHead,
	SlotNeck,
	SlotBody,
	SlotAboutBody,
	SlotLegs,
	SlotFeet,
	SlotArms,
	SlotWrist,
	SlotHands,
	SlotFingers,
	SlotWaist,
}

func EquipmentSlots() []EquipmentSlot {
	return slices.Clone(equipmentSlots)
}

// slotRank is where each slot sits in equipmentSlots, which is declared
// weapon-first and then head-to-toe on purpose: it is the order a player sees
// their own equipment, the order `role` lists what argues for what, and the
// order gear is written to a save. Sorting the slots themselves would do it
// alphabetically now that they are strings, which puts about_body first and
// wield second to last.
var slotRank = func() map[EquipmentSlot]int {
	rank := make(map[EquipmentSlot]int, len(equipmentSlots))
	for i, s := range equipmentSlots {
		rank[s] = i
	}
	return rank
}()

// Rank is the slot's place in that order. A slot nobody declared sorts last
// rather than first, so a content typo doesn't take over the top of the list.
func (e EquipmentSlot) Rank() int {
	if i, known := slotRank[e]; known {
		return i
	}
	return len(equipmentSlots)
}

// CompareSlots orders slots the way a player expects to read them, for
// slices.SortFunc.
func CompareSlots(a, b EquipmentSlot) int {
	return cmp.Compare(a.Rank(), b.Rank())
}

func ParseEquipmentSlot(s string) (EquipmentSlot, error) {
	slot := EquipmentSlot(strings.ToLower(s))
	if !slices.Contains(equipmentSlots, slot) {
		return SlotNone, fmt.Errorf("%w: %q", ErrUnknownEquipmentSlot, s)
	}
	return slot, nil
}
