package player

import (
	"github.com/stretchr/testify/suite"
	"github.com/trasa/watchmud-message/slot"
	"github.com/trasa/watchmud/object"
	"testing"
)

type InventoryTestSuite struct {
	suite.Suite
	inv        *Inventory
	definition *object.Definition
	instance   *object.Instance
}

func TestInventoryTestSuite(t *testing.T) {
	suite.Run(t, new(InventoryTestSuite))
}

func (s *InventoryTestSuite) SetupTest() {
	s.inv = NewInventory()
	s.definition = object.NewDefinition("defn", "defnName", "zone", object.Armor, []string{"alias"}, "desc", "descground", slot.Arms)
	var err error
	s.instance, err = object.NewInstance(s.definition)
	if err != nil {
		s.Assert().Fail("Failed to create instance - %s", err)
	}
}

func (s *InventoryTestSuite) Test_IsDirty() {
	s.Assert().False(s.inv.IsDirty())

	s.inv.Add(s.instance)
	s.Assert().True(s.inv.IsDirty())

	s.inv.ResetDirtyFlag()
	s.Assert().False(s.inv.IsDirty())

	s.inv.Remove(s.instance)
	s.Assert().True(s.inv.IsDirty())
}

func (s *InventoryTestSuite) Test_GetByName() {
	_ = s.inv.Add(s.instance)

	s.Assert().Equal(s.instance, s.inv.GetByNameOrAlias(s.instance.Definition.Name)[0])
}

func (s *InventoryTestSuite) Test_GetByAlias() {
	_ = s.inv.Add(s.instance)

	s.Assert().Equal(s.instance, s.inv.GetByNameOrAlias(s.instance.Definition.Aliases[0])[0])
}

func (s *InventoryTestSuite) Test_Get_NotFound() {
	_ = s.inv.Add(s.instance)

	s.Assert().Empty(s.inv.GetByNameOrAlias("doesnt_exist"))
}
