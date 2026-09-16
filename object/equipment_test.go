package object

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/trasa/watchmud/rules"
)

type EquipmentSuite struct {
	suite.Suite
	eq *Equipment
}

func TestEquipmentSuite(t *testing.T) {
	suite.Run(t, new(EquipmentSuite))
}

func (s *EquipmentSuite) SetupTest() {
	s.eq = NewEquipment()
}

func (s *EquipmentSuite) TestSlotEquipped() {
	s.Assert().False(s.eq.Equipped(rules.SlotHead))
	i := Instance{}

	s.eq.Equip(rules.SlotHead, &i)
	s.Assert().True(s.eq.Equipped(rules.SlotHead))
}

func (s *EquipmentSuite) TestSlotUnequipped() {
	i := Instance{}
	s.eq.Equip(rules.SlotHead, &i)
	s.Assert().True(s.eq.Equipped(rules.SlotHead))
	s.eq.Unequip(rules.SlotHead)
	s.Assert().False(s.eq.Equipped(rules.SlotHead))
}

func (s *EquipmentSuite) TestItemEquipped() {
	i := Instance{}
	s.eq.Equip(rules.SlotHead, &i)

	s.Assert().True(s.eq.Equipped(rules.SlotHead))
	s.Assert().True(s.eq.ItemEquipped(&i))
}
