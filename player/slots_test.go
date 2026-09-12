package player

import (
	"testing"
	"uuid"

	"github.com/stretchr/testify/suite"
	"github.com/trasa/watchmud/object"
	"github.com/trasa/watchmud/slot"
	"github.com/trasa/watchmud/thing"
)

type SlotsSuite struct {
	suite.Suite
	slots         *Slots
	slotInventory *SlotInventory
	weaponInst    *object.Instance
	armorInst     *object.Instance
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
	// hack: working around broken slots (messages) which is going away
	suite.weaponInst = &object.Instance{
		Id:         uuid.New(),
		Definition: object.NewDefinition("weapon", "weapon", "zone", object.Weapon, []string{}, "weapon", "weapon", slot.Wield),
	}
	suite.slotInventory.m.Add(suite.weaponInst)

	suite.armorInst = &object.Instance{
		Id:         uuid.New(),
		Definition: object.NewDefinition("armor", "armor", "zone", object.Armor, []string{}, "armor", "armor", slot.Head),
	}
}

func (suite *SlotsSuite) TestCantEquipYouDontHaveOne() {
	youdonthaveoneInst := &object.Instance{
		Id:         uuid.New(),
		Definition: object.NewDefinition("nothere", "nothere", "zone", object.Weapon, []string{}, "youdonthaveone", "youdonthaveone", slot.Wield),
	}
	suite.slots.Set(slot.Wield, youdonthaveoneInst)
}

func (suite *SlotsSuite) TestNotEquippableWeapon() {
	cantequipthat := &object.Instance{
		Id:         uuid.New(),
		Definition: object.NewDefinition("treasure", "treasure", "zone", object.Treasure, []string{}, "treasure", "treasure", slot.None),
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
