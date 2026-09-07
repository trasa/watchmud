package player

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/trasa/watchmud-message/slot"
	"github.com/trasa/watchmud/object"
)

type inventorySuite struct {
	suite.Suite
	inv        *Inventory
	definition *object.Definition
	instance   *object.Instance
}

func TestInventorySuite(t *testing.T) {
	suite.Run(t, new(inventorySuite))
}

func (s *inventorySuite) SetupTest() {
	s.inv = NewInventory()
	s.definition = object.NewDefinition("defn", "defnName", "zone", object.Armor, []string{"alias"}, "desc", "descground", slot.Arms)
	s.instance = object.NewInstance(s.definition)
	s.inv.Add(s.instance)
}

func (s *inventorySuite) Test_GetByName() {
	s.Assert().Equal(s.instance, s.inv.GetByNameOrAlias(s.instance.Definition.Name)[0])
}

func (s *inventorySuite) Test_GetByAlias() {
	s.Assert().Equal(s.instance, s.inv.GetByNameOrAlias(s.instance.Definition.Aliases[0])[0])
}

func (s *inventorySuite) Test_Get_NotFound() {
	s.Assert().Empty(s.inv.GetByNameOrAlias("doesnt_exist"))
}
