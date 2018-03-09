package world

import (
	"github.com/stretchr/testify/suite"
	"github.com/trasa/watchmud/client"
	"github.com/trasa/watchmud/gameserver"
	"github.com/trasa/watchmud/message"
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
	suite.world = newTestWorld()
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
	getresp := suite.player.GetSentResponse(0).(message.GetResponse)
	suite.Assert().True(getresp.Success)

	dropresp := suite.player.GetSentResponse(1).(message.DropResponse)
	suite.Assert().True(dropresp.Success)

	// player now has zero items, room has its starting two
	suite.Assert().Equal(0, len(suite.player.GetAllInventory()))
	suite.Assert().Equal(2, len(suite.world.startRoom.GetAllInventory()))
}

func (suite *HandleDropSuite) TestDrop_NoTarget() {
	// drop
	dropGameMessage, err := message.NewGameMessage(message.DropRequest{Target: ""})
	suite.Assert().NoError(err)

	dropHP := gameserver.NewHandlerParameter(suite.testClient, dropGameMessage)
	suite.world.handleDrop(dropHP)

	suite.Assert().Equal(1, suite.player.SentMessageCount())

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
	dropresp := suite.player.GetSentResponse(0).(message.DropResponse)
	suite.Assert().False(dropresp.Success)
	suite.Assert().Equal("TARGET_NOT_FOUND", dropresp.GetResultCode())
	suite.Assert().Equal(0, len(suite.player.GetAllInventory()))
}
