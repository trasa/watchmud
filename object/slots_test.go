package object

import (
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/suite"
	"github.com/trasa/watchmud/slot"
	"github.com/trasa/watchmud/thing"
	"testing"
)

type SlotsSuite struct {
	suite.Suite
	slots         *Slots
	slotInventory *SlotInventory
	weaponInst    *Instance
	armorInst     *Instance
}

type SlotInventory struct {
	m thing.Map
}

func (s *SlotInventory) Inventory() thing.Map {
	return s.m
}

func TestSlotsSuite(t *testing.T) {
	suite.Run(t, new(SlotsSuite))
}

func (suite *SlotsSuite) SetupTest() {
	suite.slotInventory = &SlotInventory{m: make(thing.Map)}
	suite.slots = NewSlots()
	suite.weaponInst = &Instance{
		InstanceId: uuid.NewV4(),
		Definition: NewDefinition("weapon", "weapon", "zone", Weapon, []string{}, "weapon", "weapon", slot.Wield),
	}
	suite.slotInventory.m.Add(suite.weaponInst)

	suite.armorInst = &Instance{
		InstanceId: uuid.NewV4(),
		Definition: NewDefinition("armor", "armor", "zone", Armor, []string{}, "armor", "armor", slot.Head),
	}
}

func (suite *SlotsSuite) TestCantEquipYouDontHaveOne() {
	youdonthaveoneInst := &Instance{
		InstanceId: uuid.NewV4(),
		Definition: NewDefinition("nothere", "nothere", "zone", Weapon, []string{}, "youdonthaveone", "youdonthaveone", slot.Wield),
	}
	suite.slots.Set(slot.Wield, youdonthaveoneInst)
}

func (suite *SlotsSuite) TestNotEquippableWeapon() {
	cantequipthat := &Instance{
		InstanceId: uuid.NewV4(),
		Definition: NewDefinition("treasure", "treasure", "zone", Treasure, []string{}, "treasure", "treasure", slot.None),
	}
	suite.slotInventory.m.Add(cantequipthat)

	// that isn't a weapon
	suite.slots.Set(slot.Wield, cantequipthat)
}

func (suite *SlotsSuite) TestSlotInUse() {
	suite.slots.Set(slot.Head, suite.armorInst)
	suite.Assert().True(suite.slots.IsSlotInUse(slot.Head))

	suite.Assert().True(suite.slots.IsItemInUse(suite.armorInst))
	suite.Assert().False(suite.slots.IsItemInUse(suite.weaponInst))
}
