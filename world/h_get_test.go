package world

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/trasa/watchmud-message"
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
	s.Assert().Equal(2, len(s.w.StartRoom.GetAllInventory()))
	s.Assert().Equal(0, len(s.p.Inventory().GetAll()))

	s.w.handleGet(s.handlerParameter(message.GetRequest{Target: "knife"}))

	s.Assert().Equal(1, len(s.r.Sent))
	resp := sent[message.GetResponse](s.T(), s.r, 0)
	s.Assert().True(resp.Success)

	// player has one item
	s.Assert().Equal(1, len(s.p.Inventory().GetAll()))
	found := s.p.Inventory().GetByNameOrAlias("knife")
	s.Assert().True(len(found) > 0)
	s.Assert().Equal("knife", found[0].Definition.Name)

	// there's one other item in the room now
	s.Assert().Equal(1, len(s.w.StartRoom.GetAllInventory()))
}

func (s *HandleGetSuite) TestAliasTarget() {
	s.w.handleGet(s.handlerParameter(message.GetRequest{Target: "iron"}))

	s.Assert().Equal(1, len(s.r.Sent))
	response := sent[message.GetResponse](s.T(), s.r, 0)
	s.Assert().True(response.Success)
	s.Assert().Equal(1, len(s.p.Inventory().GetAll()))

	found := s.p.Inventory().GetByNameOrAlias("iron_helmet")
	s.Assert().True(len(found) > 0)
	s.Assert().Equal("iron_helmet", found[0].Definition.Name)
	s.Assert().Equal(1, len(s.w.StartRoom.GetAllInventory()))
}

func (s *HandleGetSuite) TestTargetNotInRoom() {
	s.w.handleGet(s.handlerParameter(message.GetRequest{Target: "bag_of_coins"}))

	s.Assert().Equal(1, len(s.r.Sent))
	response := sent[message.GetResponse](s.T(), s.r, 0)
	s.Assert().False(response.Success)
	s.Assert().Equal("TARGET_NOT_FOUND", response.GetResultCode())

	// player has zero items still
	s.Assert().Equal(0, len(s.p.Inventory().GetAll()))

	// still two items in start room
	s.Assert().Equal(2, len(s.w.StartRoom.GetAllInventory()))
}

func (s *HandleGetSuite) TestNoTarget() {
	s.w.handleGet(s.handlerParameter(message.GetRequest{Target: ""}))

	s.Assert().Equal(1, len(s.r.Sent))
	response := sent[message.GetResponse](s.T(), s.r, 0)
	s.Assert().False(response.Success)
	s.Assert().Equal("NO_TARGET", response.GetResultCode())

	// player has zero items, start room still has 2
	s.Assert().Equal(0, len(s.p.Inventory().GetAll()))
	s.Assert().Equal(2, len(s.w.StartRoom.GetAllInventory()))
}

func (s *HandleGetSuite) TestPlayerAddFail() {
	// TODO: some sort of world-wide list of inventory definitions
	// give the player a knife to start with
	// note that two different objects should not have the same instance id
	// -- this is an arbitrary case to make the test work...
	inv, exists := s.w.StartRoom.GetInventoryByName("knife")
	s.Assert().True(exists)

	s.Assert().NoError(s.p.Inventory().Add(inv))

	s.w.handleGet(s.handlerParameter(message.GetRequest{Target: "knife"}))

	s.Assert().Equal(1, len(s.r.Sent))
	response := sent[message.GetResponse](s.T(), s.r, 0)
	s.Assert().False(response.Success)
	s.Assert().Equal("ADD_INVENTORY_ERROR", response.GetResultCode())

	// player just has one item we added at beginning of this method
	s.Assert().Equal(1, len(s.p.Inventory().GetAll()))
	// room still has its two items
	s.Assert().Equal(2, len(s.w.StartRoom.GetAllInventory()))
}

// TODO: test case for when room.Inventory.Remove fails
// need to figure out how to mock the room
