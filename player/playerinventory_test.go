package player

import (
	"github.com/stretchr/testify/suite"
	"github.com/trasa/watchmud-message/slot"
	"github.com/trasa/watchmud/object"
	"testing"
)

type PlayerInventoryTestSuite struct {
	suite.Suite
	inv        *PlayerInventory
	definition *object.Definition
	instance   *object.Instance
}

func TestPlayerInventoryTestSuite(t *testing.T) {
	suite.Run(t, new(PlayerInventoryTestSuite))
}

func (s *PlayerInventoryTestSuite) SetupTest() {
	s.inv = NewPlayerInventory()
	s.definition = object.NewDefinition("defn", "defnName", "zone", object.Armor, []string{"foo"}, "desc", "descground", slot.Arms)
	var err error
	s.instance, err = object.NewInstance(s.definition)
	if err != nil {
		s.Assert().Fail("Failed to create instance - %s", err)
	}
}

func (s *PlayerInventoryTestSuite) Test_IsDirty() {
	s.Assert().False(s.inv.IsDirty())

	s.inv.Add(s.instance)
	s.Assert().True(s.inv.IsDirty())

	s.inv.ResetDirtyFlag()
	s.Assert().False(s.inv.IsDirty())

	s.inv.Remove(s.instance)
	s.Assert().True(s.inv.IsDirty())
}
