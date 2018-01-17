package equip

import (
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/trasa/watchmud/object"
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
	suite.slots = NewSlots()
	suite.slotInventory = &SlotInventory{m: make(thing.Map)}
	suite.slots.inv = suite.slotInventory
}

func (suite *SlotsSuite) TestSetEquippedPrimaryWeapon() {
	weaponInst := &object.Instance{
		InstanceId: uuid.NewV4(),
		Definition: object.NewDefinition("weapon", "weapon", "zone", object.WEAPON, []string{}, "weapon", "weapon"),
	}
	suite.slotInventory.m.Add(weaponInst)
	err := suite.slots.Set(PrimaryHand, weaponInst)
	assert.NoError(suite.T(), err)
}

func (suite *SlotsSuite) TestSetEquippedSecondaryWeapon() {
	weaponInst := &object.Instance{
		InstanceId: uuid.NewV4(),
		Definition: object.NewDefinition("weapon", "weapon", "zone", object.WEAPON, []string{}, "weapon", "weapon"),
	}
	suite.slotInventory.m.Add(weaponInst)
	err := suite.slots.Set(SecondaryHand, weaponInst)
	assert.NoError(suite.T(), err)
}

func (suite *SlotsSuite) TestCantEquipYouDontHaveOne() {
	youdonthaveoneInst := &object.Instance{
		InstanceId: uuid.NewV4(),
		Definition: object.NewDefinition("nothere", "nothere", "zone", object.WEAPON, []string{}, "youdonthaveone", "youdonthaveone"),
	}
	assert.IsType(suite.T(), &object.InstanceNotFoundError{}, suite.slots.Set(PrimaryHand, youdonthaveoneInst))
}

func (suite *SlotsSuite) TestNotEquipableWeapon() {
	cantequipthat := &object.Instance{
		InstanceId: uuid.NewV4(),
		Definition: object.NewDefinition("treasure", "treasure", "zone", object.TREASURE, []string{}, "treasure", "treasure"),
	}
	suite.slotInventory.m.Add(cantequipthat)

	// that isn't a weapon
	assert.IsType(suite.T(), &object.InstanceNotWeaponError{}, suite.slots.Set(PrimaryHand, cantequipthat))
	assert.IsType(suite.T(), &object.InstanceNotWeaponError{}, suite.slots.Set(SecondaryHand, cantequipthat))
}
