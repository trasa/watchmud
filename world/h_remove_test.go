package world

import (
	"testing"
	"uuid"

	"github.com/stretchr/testify/suite"
	"github.com/trasa/watchmud/command"
	"github.com/trasa/watchmud/event"
	"github.com/trasa/watchmud/object"
	"github.com/trasa/watchmud/slot"
)

type RemoveTestSuite struct {
	worldTestSuite
}

func TestRemoveTestSuite(t *testing.T) {
	suite.Run(t, new(RemoveTestSuite))
}

func (s *RemoveTestSuite) helmet() *object.Instance {
	d := object.NewDefinition("iron_helmet", "iron helmet", "start", object.Armor,
		[]string{"helm"}, "iron helmet", "an iron helmet is here.", slot.Head)
	d.RoleWeights = map[string]int{"tank": 2}
	return object.NewInstance(uuid.New(), d)
}

func (s *RemoveTestSuite) TestRemoveWithNoTarget() {
	cmd := command.Remove{}
	s.w.handleRemove(s.handlerParameter(cmd), cmd)
	s.Assert().Equal(event.NoTarget, sent[event.Failed](s.T(), s.r, 0).Code)
}

func (s *RemoveTestSuite) TestRemoveSomethingNotEquipped() {
	cmd := command.Remove{Target: "iron helmet"}
	s.w.handleRemove(s.handlerParameter(cmd), cmd)
	s.Assert().Equal(event.TargetNotFound, sent[event.Failed](s.T(), s.r, 0).Code)
}

// Carrying something is not using it: remove searches the slots.
func (s *RemoveTestSuite) TestCarryingIsNotUsing() {
	s.p.Inventory().Add(s.helmet())
	cmd := command.Remove{Target: "iron helmet"}
	s.w.handleRemove(s.handlerParameter(cmd), cmd)
	s.Assert().Equal(event.TargetNotFound, sent[event.Failed](s.T(), s.r, 0).Code)
}

func (s *RemoveTestSuite) TestRemoveByAlias() {
	s.p.Slots().Set(slot.Head, s.helmet())
	cmd := command.Remove{Target: "helm"}
	s.w.handleRemove(s.handlerParameter(cmd), cmd)
	s.Assert().Equal("iron helmet", sent[event.Removed](s.T(), s.r, 0).Item)
	s.Assert().Nil(s.p.Slots().Get(slot.Head))
}

// The slot has to be genuinely free afterwards, not holding a nil that still
// reads as occupied -- otherwise you could take a helmet off once and never
// wear anything on your head again.
func (s *RemoveTestSuite) TestTheSlotIsUsableAgain() {
	s.p.Slots().Set(slot.Head, s.helmet())
	cmd := command.Remove{Target: "iron helmet"}
	s.w.handleRemove(s.handlerParameter(cmd), cmd)
	s.Require().False(s.p.Slots().IsSlotInUse(slot.Head))

	wear := command.Wear{Target: "helm"}
	s.p.Inventory().Add(s.helmet())
	s.w.handleWear(s.handlerParameter(wear), wear)
	sent[event.Worn](s.T(), s.r, 1)
	s.Assert().True(s.p.Slots().IsSlotInUse(slot.Head))
}

// Removing the gear removes the role: this is the whole reason remove exists.
func (s *RemoveTestSuite) TestRemovingTheGearRemovesTheRole() {
	s.p.Slots().Set(slot.Head, s.helmet())
	s.Assert().Equal("Tank", s.w.roleName(s.p.RoleWeights()))

	cmd := command.Remove{Target: "iron helmet"}
	s.w.handleRemove(s.handlerParameter(cmd), cmd)
	s.Assert().Equal("", s.w.roleName(s.p.RoleWeights()))
}

// And the equipment listing doesn't trip over the emptied slot.
func (s *RemoveTestSuite) TestEquipmentListingAfterRemove() {
	s.p.Slots().Set(slot.Head, s.helmet())
	cmd := command.Remove{Target: "iron helmet"}
	s.w.handleRemove(s.handlerParameter(cmd), cmd)

	s.w.handleShowEquipment(s.handlerParameter(command.ShowEquipment{}), command.ShowEquipment{})
	s.Assert().Empty(sent[event.Equipment](s.T(), s.r, 1).Items)
}
