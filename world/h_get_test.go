package world

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/trasa/watchmud-message"
	"github.com/trasa/watchmud/gameserver"
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

func (s *HandleGetSuite) handlerParameter(target string) *gameserver.HandlerParameter {
	return s.worldTestSuite.handlerParameter(message.GetRequest{Target: target})
}

func (s *HandleGetSuite) TestGet_success() {
	// start off with two items in the room and zero in the player
	s.Assert().Equal(2, len(s.w.StartRoom.GetAllInventory()))
	s.Assert().Equal(0, len(s.p.Inventory().GetAll()))

	getHP := s.handlerParameter("knife")
	s.w.handleGet(getHP)

	s.Assert().Equal(1, len(s.r.Sent))
	resp := s.r.Sent[0].(message.GetResponse)
	s.Assert().True(resp.Success)

	// player has one item
	s.Assert().Equal(1, len(s.p.Inventory().GetAll()))
	foundinv := s.p.Inventory().GetByNameOrAlias("knife")
	s.Assert().True(len(foundinv) > 0)
	s.Assert().Equal("knife", foundinv[0].Definition.Name)

	// there's one other item in the room now
	s.Assert().Equal(1, len(s.w.StartRoom.GetAllInventory()))
}

func (s *HandleGetSuite) TestGet_AliasTarget() {
	getHP := s.handlerParameter("iron")
	s.w.handleGet(getHP)

	s.Assert().Equal(1, len(s.r.Sent))
	resp := s.r.Sent[0].(message.GetResponse)
	s.Assert().True(resp.Success)
	s.Assert().Equal(1, len(s.p.Inventory().GetAll()))

	foundinv := s.p.Inventory().GetByNameOrAlias("iron_helmet")
	s.Assert().True(len(foundinv) > 0)
	s.Assert().Equal("iron_helmet", foundinv[0].Definition.Name)
	s.Assert().Equal(1, len(s.w.StartRoom.GetAllInventory()))
}

func (s *HandleGetSuite) TestGet_targetNotInRoom() {
	getHP := s.handlerParameter("bag_of_coins")
	s.w.handleGet(getHP)

	s.Assert().Equal(1, len(s.r.Sent))
	resp := s.r.Sent[0].(message.GetResponse)
	s.Assert().True(resp.Success)
	s.Assert().Equal("TARGET_NOT_FOUND", resp.GetResultCode())

	// player has zero items still
	s.Assert().Equal(0, len(s.p.Inventory().GetAll()))

	// still two items in start room
	s.Assert().Equal(2, len(s.w.StartRoom.GetAllInventory()))
}

func (s *HandleGetSuite) TestGet_NoTarget() {
	getHP := s.handlerParameter("")
	s.w.handleGet(getHP)

	s.Assert().Equal(1, len(s.r.Sent))
	resp := s.r.Sent[0].(message.GetResponse)
	s.Assert().True(resp.Success)
	s.Assert().Equal("NO_TARGET", resp.GetResultCode())

	// player has zero items, start room still has 2
	s.Assert().Equal(0, len(s.p.Inventory().GetAll()))
	s.Assert().Equal(2, len(s.w.StartRoom.GetAllInventory()))
}

func (s *HandleGetSuite) TestGet_PlayerAddFail() {
	// TODO: some sort of world-wide list of inventory definitions
	// give the player a knife to start with
	// note that two different objects should not have the same instance id
	// -- this is an arbitrary case to make the test work...
	inv, exists := s.w.StartRoom.GetInventoryByName("knife")
	s.Assert().True(exists)

	s.Assert().NoError(s.p.Inventory().Add(inv))

	getHP := s.handlerParameter("knife")
	s.w.handleGet(getHP)

	s.Assert().Equal(1, len(s.r.Sent))
	resp := s.r.Sent[0].(message.GetResponse)
	s.Assert().False(resp.Success)
	s.Assert().Equal("ADD_INVENTORY_ERROR", resp.GetResultCode())

	// player just has one item we added at beginning of this method
	s.Assert().Equal(1, len(s.p.Inventory().GetAll()))
	// room still has its two items
	s.Assert().Equal(2, len(s.w.StartRoom.GetAllInventory()))
}

// TODO: test case for when room.Inventory.Remove fails
// need to figure out how to mock the room
