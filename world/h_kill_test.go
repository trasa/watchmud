package world

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/trasa/watchmud-message"
	"github.com/trasa/watchmud/client"
	"github.com/trasa/watchmud/gameserver"
	"github.com/trasa/watchmud/mobile"
	"github.com/trasa/watchmud/player"
	"github.com/trasa/watchmud/spaces"
	"testing"
)

type HandleKillSuite struct {
	suite.Suite
	world      *World
	player     *player.TestPlayer
	testClient *client.TestClient
}

func TestHandleKillSuite(t *testing.T) {
	suite.Run(t, new(HandleKillSuite))
}

func (suite *HandleKillSuite) SetupTest() {
	suite.world, _ = newTestWorld()
	suite.player = player.NewTestPlayer("testdood")
	suite.world.AddPlayer(suite.player)
	suite.testClient = client.NewTestClient(suite.player)
}

func newKillRequestHandleParameter(t *testing.T, c *client.TestClient, target string) *gameserver.HandlerParameter {
	msg, err := message.NewGameMessage(message.KillRequest{Target: target})
	assert.NoError(t, err)
	return gameserver.NewHandlerParameter(c, msg)
}

func (suite *HandleKillSuite) TestSuccess() {
	mob, _ := suite.world.StartRoom.FindMobile("target")
	killHP := newKillRequestHandleParameter(suite.T(), suite.testClient, "target")

	suite.world.handleKill(killHP)

	suite.Assert().Equal(1, suite.player.SentMessageCount())
	resp := suite.player.GetSentResponse(0).(message.KillResponse)
	suite.Assert().True(resp.Success)
	suite.Assert().Equal("OK", resp.ResultCode)

	suite.Assert().True(suite.world.fightLedger.IsFighting(suite.player))
	suite.Assert().Equal(mob, suite.world.fightLedger.GetFight(suite.player).Fightee)
	suite.Assert().Equal(suite.player, suite.world.fightLedger.GetFight(suite.player).Fighter)
}

func (suite *HandleKillSuite) TestAlreadyFighting() {
	mob, _ := suite.world.StartRoom.FindMobile("target")
	suite.world.fightLedger.Fight(suite.player, mob, suite.world.StartRoom.Zone.Id, suite.world.StartRoom.Id)

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

func (suite *HandleKillSuite) TestNoFight() {
	mob, _ := suite.world.StartRoom.FindMobile("target")
	mob.Definition.SetFlag(mobile.FlagNoFight)
	killHP := newKillRequestHandleParameter(suite.T(), suite.testClient, "target")

	suite.world.handleKill(killHP)

	suite.Assert().Equal(1, suite.player.SentMessageCount())
	resp := suite.player.GetSentResponse(0).(message.KillResponse)

	suite.Assert().False(resp.Success)
	suite.Assert().Equal("NO_FIGHT", resp.ResultCode)
}

func (suite *HandleKillSuite) TestNoFightInRoom() {

	suite.world.StartRoom.SetFlag(spaces.RoomFlagNoFight)

	killHP := newKillRequestHandleParameter(suite.T(), suite.testClient, "target")

	suite.world.handleKill(killHP)

	suite.Assert().Equal(1, suite.player.SentMessageCount())
	resp := suite.player.GetSentResponse(0).(message.KillResponse)

	suite.Assert().False(resp.Success)
	suite.Assert().Equal("NO_FIGHT_ROOM", resp.ResultCode)
}
