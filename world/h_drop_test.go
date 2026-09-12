package world

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/trasa/watchmud/command"
	"github.com/trasa/watchmud/event"
	"github.com/trasa/watchmud/slot"
)

type HandleDropSuite struct {
	worldTestSuite
}

func TestHandleDropSuite(t *testing.T) {
	suite.Run(t, new(HandleDropSuite))
}

func (s *HandleDropSuite) SetupTest() {
	s.worldTestSuite.SetupTest()
}

// get picks something up so there is something to drop. Still on the legacy
// path until handleGet is converted.
func (s *HandleDropSuite) get(target string) {
	s.T().Helper()
	cmd := command.Get{Target: target}
	s.w.handleGet(s.handlerParameter(cmd), cmd)
}

func (s *HandleDropSuite) drop(target string) {
	s.T().Helper()
	cmd := command.Drop{Target: target}
	s.w.handleDrop(s.handlerParameter(cmd), cmd)
}

func (s *HandleDropSuite) TestSuccess() {
	s.get("knife")
	s.drop("knife")

	dropped := sent[event.Dropped](s.T(), s.r, 1)
	s.Assert().Equal("testdood", dropped.Actor)
	s.Assert().Equal("knife", dropped.Item)

	// player now has zero items, room has its starting two
	s.Assert().Equal(0, len(s.p.Inventory().GetAll()))
	s.Assert().Equal(2, len(s.w.StartRoom.Inventory.GetAll()))
}

func (s *HandleDropSuite) TestAlias() {
	s.get("helmet")
	s.drop("helmet")

	s.Assert().Equal(2, len(s.r.Sent))
	sent[event.Dropped](s.T(), s.r, 1)

	// player now has zero items, room has its starting two
	s.Assert().Equal(0, len(s.p.Inventory().GetAll()))
	s.Assert().Equal(2, len(s.w.StartRoom.Inventory.GetAll()))
}

func (s *HandleDropSuite) TestNoTarget() {
	s.drop("")

	s.Assert().Equal(1, len(s.r.Sent))
	failed := sent[event.Failed](s.T(), s.r, 0)
	s.Assert().Equal("drop", failed.Verb)
	s.Assert().Equal(event.NoTarget, failed.Code)
}

func (s *HandleDropSuite) TestNotFound() {
	// drop (but you don't have one)
	s.drop("knife")

	s.Assert().Equal(1, len(s.r.Sent))
	failed := sent[event.Failed](s.T(), s.r, 0)
	s.Assert().Equal("drop", failed.Verb)
	s.Assert().Equal(event.TargetNotFound, failed.Code)
	s.Assert().Equal(0, len(s.p.Inventory().GetAll()))
}

func (s *HandleDropSuite) TestInUse() {
	s.get("knife")

	// now wield the knife
	equip := command.Equip{Target: "knife", Slot: slot.Wield}
	s.w.handleEquip(s.handlerParameter(equip), equip)

	s.drop("knife")

	failed := sent[event.Failed](s.T(), s.r, 2)
	s.Assert().Equal(event.TargetInUse, failed.Code)
	s.Assert().Equal(1, len(s.p.Inventory().GetAll()))
}

func (s *HandleDropSuite) TestInUseMultipleItems() {
	// drop knife
	// - there are two knives on you,
	//   one you are holding, another in your inventory list.
	// - drop the one that is just in the inventory list and not in use
	//suite.Assert().Fail("TODO implement me")
}
