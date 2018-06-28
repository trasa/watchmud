package world

import (
	"github.com/stretchr/testify/suite"
	"github.com/trasa/watchmud/player"
	"github.com/trasa/watchmud/client"
	"testing"
	"github.com/trasa/watchmud/gameserver"
	"github.com/trasa/watchmud/message"
	"github.com/stretchr/testify/assert"
)

type HandleKillSuite struct {
	suite.Suite
	world *World
	player *player.TestPlayer
	testClient *client.TestClient
}

func TestHandleKillSuite(t *testing.T) {
	suite.Run(t, new(HandleKillSuite))
}

func (suite *HandleKillSuite) SetupTest() {
	suite.world = newTestWorld()
	suite.player = player.NewTestPlayer("testdood")
	suite.world.AddPlayer(suite.player)
	suite.testClient = client.NewTestClient(suite.player)
}

func newKillRequestHandleParameter(t *testing.T, c *client.TestClient, target string) *gameserver.HandlerParameter {
	msg, err := message.NewGameMessage(message.KillRequest{Target: target})
	assert.NoError(t, err)
	return gameserver.NewHandlerParameter(c, msg)
}

func (suite *HandleKillSuite) TestAlreadyFighting() {
	suite.player.SetFighting(true)
	killHP := newKillRequestHandleParameter(suite.T(), suite.testClient, "targetMob")

	suite.world.handleKill(killHP)

	suite.Assert().Equal(1, suite.player.SentMessageCount())
	resp := suite.player.GetSentResponse(0).(message.KillResponse)
	suite.Assert().False(resp.Success)
	suite.Assert().Equal("ALREADY_FIGHTING", resp.ResultCode)
}

func (suite *HandleKillSuite) TestNoTarget() {
	killHP := newKillRequestHandleParameter(suite.T(), suite.testClient, "targetMob")

	suite.world.handleKill(killHP)

	suite.Assert().Equal(1, suite.player.SentMessageCount())
	resp := suite.player.GetSentResponse(0).(message.KillResponse)
	suite.Assert().False(resp.Success)
	suite.Assert().Equal("TARGET_NOT_FOUND", resp.ResultCode)
}