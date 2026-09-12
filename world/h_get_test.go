package world

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/trasa/watchmud/command"
	"github.com/trasa/watchmud/event"
)

type HandleGetSuite struct {
	worldTestSuite
}

func TestHandleGetSuite(t *testing.T) {
	suite.Run(t, new(HandleGetSuite))
}

func (s *HandleGetSuite) SetupTest() {
	s.worldTestSuite.SetupTest()
}

func (s *HandleGetSuite) TestSuccess() {
	// start off with two items in the room and zero in the player
	s.Assert().Equal(2, len(s.w.StartRoom.Inventory.GetAll()))
	s.Assert().Equal(0, len(s.p.Inventory().GetAll()))

	cmd := command.Get{Target: "knife"}
	s.w.handleGet(s.handlerParameter(cmd), cmd)

	s.Assert().Equal(1, len(s.r.Sent))
	got := sent[event.Got](s.T(), s.r, 0)
	s.Assert().Equal("testdood", got.Actor)

	// player has one item
	s.Assert().Equal(1, len(s.p.Inventory().GetAll()))
	found := s.p.Inventory().GetByNameOrAlias("knife")
	s.Assert().True(len(found) > 0)
	s.Assert().Equal("knife", found[0].Definition.Name)

	// there's one other item in the room now
	s.Assert().Equal(1, len(s.w.StartRoom.Inventory.GetAll()))
}

func (s *HandleGetSuite) TestAliasTarget() {
	cmd := command.Get{Target: "iron"}
	s.w.handleGet(s.handlerParameter(cmd), cmd)

	s.Assert().Equal(1, len(s.r.Sent))
	sent[event.Got](s.T(), s.r, 0)
	s.Assert().Equal(1, len(s.p.Inventory().GetAll()))

	found := s.p.Inventory().GetByNameOrAlias("helmet")
	s.Assert().True(len(found) > 0)
	s.Assert().Equal("helmet", found[0].Definition.Name)
	s.Assert().Equal(1, len(s.w.StartRoom.Inventory.GetAll()))
}

func (s *HandleGetSuite) TestTargetNotInRoom() {
	cmd := command.Get{Target: "bag_of_coins"}
	s.w.handleGet(s.handlerParameter(cmd), cmd)

	s.Assert().Equal(1, len(s.r.Sent))
	failed := sent[event.Failed](s.T(), s.r, 0)
	s.Assert().Equal("get", failed.Verb)
	s.Assert().Equal(event.TargetNotFound, failed.Code)

	// player has zero items still
	s.Assert().Equal(0, len(s.p.Inventory().GetAll()))

	// still two items in start room
	s.Assert().Equal(2, len(s.w.StartRoom.Inventory.GetAll()))
}

func (s *HandleGetSuite) TestNoTarget() {
	cmd := command.Get{Target: ""}
	s.w.handleGet(s.handlerParameter(cmd), cmd)

	s.Assert().Equal(1, len(s.r.Sent))
	failed := sent[event.Failed](s.T(), s.r, 0)
	s.Assert().Equal(event.NoTarget, failed.Code)

	// player has zero items, start room still has 2
	s.Assert().Equal(0, len(s.p.Inventory().GetAll()))
	s.Assert().Equal(2, len(s.w.StartRoom.Inventory.GetAll()))
}
