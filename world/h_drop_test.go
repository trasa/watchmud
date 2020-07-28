package world

import (
	"github.com/stretchr/testify/suite"
	"github.com/trasa/watchmud-message"
	"github.com/trasa/watchmud-message/slot"
	"github.com/trasa/watchmud/client"
	"github.com/trasa/watchmud/gameserver"
	"github.com/trasa/watchmud/player"
	"testing"
)

type HandleDropSuite struct {
	suite.Suite
	world      *World
	player     *player.TestPlayer
	testClient *client.TestClient
}

func TestHandleDropSuite(t *testing.T) {
	suite.Run(t, new(HandleDropSuite))
}

func (suite *HandleDropSuite) SetupTest() {
	suite.world, _ = newTestWorld()
	suite.player = player.NewTestPlayer("foo")
	suite.world.AddPlayer(suite.player)
	suite.testClient = client.NewTestClient(suite.player)
}

func (suite *HandleDropSuite) TestDrop_success() {
	// get first
	getGameMessage, err := message.NewGameMessage(
		message.GetRequest{
			Target: "knife",
		})
	suite.Assert().NoError(err)

	getHP := gameserver.NewHandlerParameter(suite.testClient, getGameMessage)
	suite.world.handleGet(getHP)

	// now drop
	dropGameMessage, err := message.NewGameMessage(
		message.DropRequest{
			Target: "knife",
		},
	)
	suite.Assert().NoError(err)

	dropHP := gameserver.NewHandlerParameter(suite.testClient, dropGameMessage)
	suite.world.handleDrop(dropHP)

	suite.Assert().Equal(2, suite.player.SentMessageCount())
	//noinspection GoVetCopyLock
	getresp := suite.player.GetSentResponse(0).(message.GetResponse)
	suite.Assert().True(getresp.Success)

	//noinspection GoVetCopyLock
	dropresp := suite.player.GetSentResponse(1).(message.DropResponse)
	suite.Assert().True(dropresp.Success)

	// player now has zero items, room has its starting two
	suite.Assert().Equal(0, len(suite.player.GetInventory().GetAll()))
	suite.Assert().Equal(2, len(suite.world.StartRoom.GetAllInventory()))
}

func (suite *HandleDropSuite) TestDrop_Alias() {
	// get first
	getGameMessage, err := message.NewGameMessage(
		message.GetRequest{
			Target: "helmet",
		})
	suite.Assert().NoError(err)

	getHP := gameserver.NewHandlerParameter(suite.testClient, getGameMessage)
	suite.world.handleGet(getHP)

	// now drop
	dropGameMessage, err := message.NewGameMessage(
		message.DropRequest{
			Target: "helmet",
		},
	)
	suite.Assert().NoError(err)

	dropHP := gameserver.NewHandlerParameter(suite.testClient, dropGameMessage)
	suite.world.handleDrop(dropHP)

	suite.Assert().Equal(2, suite.player.SentMessageCount())
	//noinspection GoVetCopyLock
	getresp := suite.player.GetSentResponse(0).(message.GetResponse)
	suite.Assert().True(getresp.Success)

	//noinspection GoVetCopyLock
	dropresp := suite.player.GetSentResponse(1).(message.DropResponse)
	suite.Assert().True(dropresp.Success)

	// player now has zero items, room has its starting two
	suite.Assert().Equal(0, len(suite.player.GetInventory().GetAll()))
	suite.Assert().Equal(2, len(suite.world.StartRoom.GetAllInventory()))
}

func (suite *HandleDropSuite) TestDrop_NoTarget() {
	// drop
	dropGameMessage, err := message.NewGameMessage(message.DropRequest{Target: ""})
	suite.Assert().NoError(err)

	dropHP := gameserver.NewHandlerParameter(suite.testClient, dropGameMessage)
	suite.world.handleDrop(dropHP)

	suite.Assert().Equal(1, suite.player.SentMessageCount())

	//noinspection GoVetCopyLock
	dropresp := suite.player.GetSentResponse(0).(message.DropResponse)
	suite.Assert().False(dropresp.Success)
	suite.Assert().Equal("NO_TARGET", dropresp.GetResultCode())
}

func (suite *HandleDropSuite) TestDrop_NotFound() {
	// drop (but you don't have one)
	dropGameMessage, err := message.NewGameMessage(message.DropRequest{Target: "knife"})
	suite.Assert().NoError(err)

	dropHP := gameserver.NewHandlerParameter(suite.testClient, dropGameMessage)
	suite.world.handleDrop(dropHP)

	suite.Assert().Equal(1, suite.player.SentMessageCount())
	//noinspection GoVetCopyLock
	dropresp := suite.player.GetSentResponse(0).(message.DropResponse)
	suite.Assert().False(dropresp.Success)
	suite.Assert().Equal("TARGET_NOT_FOUND", dropresp.GetResultCode())
	suite.Assert().Equal(0, len(suite.player.GetInventory().GetAll()))
}

func (suite *HandleDropSuite) TestDrop_InUse() {
	// get first
	getGameMessage, err := message.NewGameMessage(
		message.GetRequest{
			Target: "knife",
		})
	suite.Assert().NoError(err)

	getHP := gameserver.NewHandlerParameter(suite.testClient, getGameMessage)
	suite.world.handleGet(getHP)

	// now wield the knife
	wieldGameMessage, err := message.NewGameMessage(
		message.EquipRequest{
			Target:       "knife",
			SlotLocation: int32(slot.Wield),
		})
	suite.Assert().NoError(err)
	suite.world.handleEquip(gameserver.NewHandlerParameter(suite.testClient, wieldGameMessage))

	// now drop
	dropGameMessage, err := message.NewGameMessage(
		message.DropRequest{
			Target: "knife",
		},
	)
	suite.Assert().NoError(err)
	dropHP := gameserver.NewHandlerParameter(suite.testClient, dropGameMessage)

	suite.world.handleDrop(dropHP)
	//noinspection GoVetCopyLock
	dropresp := suite.player.GetSentResponse(2).(message.DropResponse)
	suite.Assert().False(dropresp.Success)
	suite.Assert().Equal("TARGET_IN_USE", dropresp.GetResultCode())
	suite.Assert().Equal(1, len(suite.player.GetInventory().GetAll()))
}

func (suite *HandleDropSuite) TestDrop_InUse_But_Multiple_Items() {
	// drop knife
	// - there are two knives on you,
	//   one you are holding, another in your inventory list.
	// - drop the one that is just in the inventory list and not in use
	//suite.Assert().Fail("TODO implement me")
}
