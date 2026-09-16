package rules

import (
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

func ParseEquipmentSlot(s string) (EquipmentSlot, error) {
	slot := EquipmentSlot(strings.ToLower(s))
	if !slices.Contains(equipmentSlots, slot) {
		return SlotNone, fmt.Errorf("%w: %q", ErrUnknownEquipmentSlot, s)
	}
	return slot, nil
}
