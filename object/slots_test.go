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
}

func (suite *SlotsSuite) TestSetEquippedPrimaryWeapon() {
	weaponInst := &Instance{
		InstanceId: uuid.Must(uuid.NewV4()),
		Definition: NewDefinition("weapon", "weapon", "zone", Weapon, []string{}, "weapon", "weapon", slot.Wield),
	}
	suite.slotInventory.m.Add(weaponInst)
	suite.slots.Set(slot.Wield, weaponInst)
}

func (suite *SlotsSuite) TestSetEquippedSecondaryWeapon() {
	weaponInst := &Instance{
		InstanceId: uuid.Must(uuid.NewV4()),
		Definition: NewDefinition("weapon", "weapon", "zone", Weapon, []string{}, "weapon", "weapon", slot.Wield),
	}
	suite.slotInventory.m.Add(weaponInst)
	suite.slots.Set(slot.Wield, weaponInst)
}

func (suite *SlotsSuite) TestCantEquipYouDontHaveOne() {
	youdonthaveoneInst := &Instance{
		InstanceId: uuid.Must(uuid.NewV4()),
		Definition: NewDefinition("nothere", "nothere", "zone", Weapon, []string{}, "youdonthaveone", "youdonthaveone", slot.Wield),
	}
	suite.slots.Set(slot.Wield, youdonthaveoneInst)
}

func (suite *SlotsSuite) TestNotEquipableWeapon() {
	cantequipthat := &Instance{
		InstanceId: uuid.Must(uuid.NewV4()),
		Definition: NewDefinition("treasure", "treasure", "zone", Treasure, []string{}, "treasure", "treasure", slot.None),
	}
	suite.slotInventory.m.Add(cantequipthat)

	// that isn't a weapon
	suite.slots.Set(slot.Wield, cantequipthat)
}
