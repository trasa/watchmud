package world

import (
	"github.com/stretchr/testify/suite"
	"github.com/trasa/watchmud-message"
	"github.com/trasa/watchmud-message/direction"
	"github.com/trasa/watchmud/client"
	"github.com/trasa/watchmud/gameserver"
	"github.com/trasa/watchmud/player"
	"testing"
)

type HandleMoveSuite struct {
	suite.Suite
	world      *World
	player     *player.TestPlayer
	testClient *client.TestClient
}

func TestHandleMoveSuite(t *testing.T) {
	suite.Run(t, new(HandleMoveSuite))
}

func (suite *HandleMoveSuite) SetupTest() {
	suite.world = newTestWorld()
	suite.player = player.NewTestPlayer("p")
	suite.world.AddPlayer(suite.player)
	suite.testClient = client.NewTestClient(suite.player)
}

func (suite *HandleMoveSuite) newMoveRequestHandlerParameter(dir direction.Direction) *gameserver.HandlerParameter {
	msg, err := message.NewGameMessage(message.MoveRequest{Direction: int32(dir)})
	suite.Assert().NoError(err)
	return gameserver.NewHandlerParameter(suite.testClient, msg)
}

func (suite *HandleMoveSuite) TestMove_butYouCant() {
	suite.world.handleMove(suite.newMoveRequestHandlerParameter(direction.NORTH))

	suite.Assert().Equal(1, suite.player.SentMessageCount())

	resp := suite.player.GetSentResponse(0).(message.MoveResponse)
	suite.Assert().False(resp.Success)
	suite.Assert().Equal(resp.ResultCode, "CANT_GO_THAT_WAY")
}

func (suite *HandleMoveSuite) TestMoveWhileFighting() {
	otherPlayer := player.NewTestPlayer("other")
	suite.world.AddPlayer(otherPlayer)
	suite.world.fightLedger.Fight(suite.player, otherPlayer, suite.world.startRoom.Zone.Id, suite.world.startRoom.Id)

	suite.world.handleMove(suite.newMoveRequestHandlerParameter(direction.NORTH))

	suite.Assert().Equal(1, suite.player.SentMessageCount())

	resp := suite.player.GetSentResponse(0).(message.MoveResponse)
	suite.Assert().False(resp.Success)
	suite.Assert().Equal(resp.ResultCode, "IN_A_FIGHT")
}
