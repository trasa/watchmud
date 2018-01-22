package object

import (
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
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
	suite.slots = NewSlots(suite.slotInventory)
}

func (suite *SlotsSuite) TestSetEquippedPrimaryWeapon() {
	weaponInst := &Instance{
		InstanceId: uuid.NewV4(),
		Definition: NewDefinition("weapon", "weapon", "zone", WEAPON, []string{}, "weapon", "weapon", slot.PrimaryHand),
	}
	suite.slotInventory.m.Add(weaponInst)
	err := suite.slots.Set(slot.PrimaryHand, weaponInst)
	assert.NoError(suite.T(), err)
}

func (suite *SlotsSuite) TestSetEquippedSecondaryWeapon() {
	weaponInst := &Instance{
		InstanceId: uuid.NewV4(),
		Definition: NewDefinition("weapon", "weapon", "zone", WEAPON, []string{}, "weapon", "weapon", slot.PrimaryHand),
	}
	suite.slotInventory.m.Add(weaponInst)
	err := suite.slots.Set(slot.SecondaryHand, weaponInst)
	assert.NoError(suite.T(), err)
}

func (suite *SlotsSuite) TestCantEquipYouDontHaveOne() {
	youdonthaveoneInst := &Instance{
		InstanceId: uuid.NewV4(),
		Definition: NewDefinition("nothere", "nothere", "zone", WEAPON, []string{}, "youdonthaveone", "youdonthaveone", slot.PrimaryHand),
	}
	assert.IsType(suite.T(), &InstanceNotFoundError{}, suite.slots.Set(slot.PrimaryHand, youdonthaveoneInst))
}

func (suite *SlotsSuite) TestNotEquipableWeapon() {
	cantequipthat := &Instance{
		InstanceId: uuid.NewV4(),
		Definition: NewDefinition("treasure", "treasure", "zone", TREASURE, []string{}, "treasure", "treasure", slot.None),
	}
	suite.slotInventory.m.Add(cantequipthat)

	// that isn't a weapon
	assert.IsType(suite.T(), &InstanceNotWeaponError{}, suite.slots.Set(slot.PrimaryHand, cantequipthat))
	assert.IsType(suite.T(), &InstanceNotWeaponError{}, suite.slots.Set(slot.SecondaryHand, cantequipthat))
}
