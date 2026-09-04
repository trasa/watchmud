package world

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/trasa/watchmud-message"
	"github.com/trasa/watchmud-message/slot"
	"github.com/trasa/watchmud/gameserver"
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

func (s *HandleDropSuite) TestSuccess() {
	// get first
	s.w.handleGet(s.handlerParameter(message.GetRequest{Target: "knife"}))

	// now drop
	s.w.handleDrop(s.handlerParameter(message.DropRequest{Target: "knife"}))

	getResponse := sent[message.GetResponse](s.T(), s.r, 0)
	s.Assert().True(getResponse.Success)

	dropResponse := sent[message.DropResponse](s.T(), s.r, 1)
	s.Assert().True(dropResponse.Success)

	// player now has zero items, room has its starting two
	s.Assert().Equal(0, len(s.p.Inventory().GetAll()))
	s.Assert().Equal(2, len(s.w.StartRoom.GetAllInventory()))
}

func (s *HandleDropSuite) TestAlias() {
	// get first
	getGameMessage, err := message.NewGameMessage(
		message.GetRequest{
			Target: "helmet",
		})
	s.Assert().NoError(err)

	getHP := gameserver.NewHandlerParameter(s.c, getGameMessage)
	s.w.handleGet(getHP)

	// now drop
	dropGameMessage, err := message.NewGameMessage(
		message.DropRequest{
			Target: "helmet",
		},
	)
	s.Assert().NoError(err)

	dropHP := gameserver.NewHandlerParameter(s.c, dropGameMessage)
	s.w.handleDrop(dropHP)

	s.Assert().Equal(2, len(s.r.Sent))
	//noinspection GoVetCopyLock
	getresp := s.r.Sent[0].(message.GetResponse)
	s.Assert().True(getresp.Success)

	//noinspection GoVetCopyLock
	dropresp := s.r.Sent[1].(message.DropResponse)
	s.Assert().True(dropresp.Success)

	// player now has zero items, room has its starting two
	s.Assert().Equal(0, len(s.p.Inventory().GetAll()))
	s.Assert().Equal(2, len(s.w.StartRoom.GetAllInventory()))
}

func (s *HandleDropSuite) TestNoTarget() {
	// drop
	dropGameMessage, err := message.NewGameMessage(message.DropRequest{Target: ""})
	s.Assert().NoError(err)

	dropHP := gameserver.NewHandlerParameter(s.c, dropGameMessage)
	s.w.handleDrop(dropHP)

	s.Assert().Equal(1, len(s.r.Sent))

	dropresp := s.r.Sent[0].(message.DropResponse)
	s.Assert().False(dropresp.Success)
	s.Assert().Equal("NO_TARGET", dropresp.GetResultCode())
}

func (s *HandleDropSuite) TestNotFound() {
	// drop (but you don't have one)
	dropGameMessage, err := message.NewGameMessage(message.DropRequest{Target: "knife"})
	s.Assert().NoError(err)

	dropHP := gameserver.NewHandlerParameter(s.c, dropGameMessage)
	s.w.handleDrop(dropHP)

	s.Assert().Equal(1, len(s.r.Sent))
	//noinspection GoVetCopyLock
	dropresp := s.r.Sent[0].(message.DropResponse)
	s.Assert().False(dropresp.Success)
	s.Assert().Equal("TARGET_NOT_FOUND", dropresp.GetResultCode())
	s.Assert().Equal(0, len(s.p.Inventory().GetAll()))
}

func (s *HandleDropSuite) TestInUse() {
	// get first
	getGameMessage, err := message.NewGameMessage(
		message.GetRequest{
			Target: "knife",
		})
	s.Assert().NoError(err)

	getHP := gameserver.NewHandlerParameter(s.c, getGameMessage)
	s.w.handleGet(getHP)

	// now wield the knife
	wieldGameMessage, err := message.NewGameMessage(
		message.EquipRequest{
			Target:       "knife",
			SlotLocation: int32(slot.Wield),
		})
	s.Assert().NoError(err)
	s.w.handleEquip(gameserver.NewHandlerParameter(s.c, wieldGameMessage))

	// now drop
	dropGameMessage, err := message.NewGameMessage(
		message.DropRequest{
			Target: "knife",
		},
	)
	s.Assert().NoError(err)
	dropHP := gameserver.NewHandlerParameter(s.c, dropGameMessage)

	s.w.handleDrop(dropHP)
	dropresp := s.r.Sent[2].(message.DropResponse)
	s.Assert().False(dropresp.Success)
	s.Assert().Equal("TARGET_IN_USE", dropresp.GetResultCode())
	s.Assert().Equal(1, len(s.p.Inventory().GetAll()))
}

func (s *HandleDropSuite) TestInUseMultipleItems() {
	// drop knife
	// - there are two knives on you,
	//   one you are holding, another in your inventory list.
	// - drop the one that is just in the inventory list and not in use
	//suite.Assert().Fail("TODO implement me")
}
