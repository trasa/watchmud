package world

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/trasa/watchmud-message"
	"github.com/trasa/watchmud/client"
	"github.com/trasa/watchmud/gameserver"
	"github.com/trasa/watchmud/player"
	"testing"
)

type HandleEquipSuite struct {
	suite.Suite
	world      *World
	player     *player.TestPlayer
	testClient *client.TestClient
	msg        *message.GameMessage
	handle     *gameserver.HandlerParameter
}

func TestHandleEquipSuite(t *testing.T) {
	suite.Run(t, new(HandleEquipSuite))
}

func (suite *HandleEquipSuite) SetupTest() {
	suite.world, _ = newTestWorld()
	suite.player = player.NewTestPlayer("foo")
	suite.world.AddPlayer(suite.player)
	suite.testClient = client.NewTestClient(suite.player)

	msg, err := message.NewGameMessage(message.EquipRequest{})
	assert.NoError(suite.T(), err)
	suite.msg = msg
	suite.handle = gameserver.NewHandlerParameter(suite.testClient, suite.msg)
}

func (suite *HandleEquipSuite) TestNoSlot() {

	suite.world.handleEquip(suite.handle)

	assert.Equal(suite.T(), 1, suite.player.SentMessageCount())
	resp := suite.player.GetSentResponse(0).(message.EquipResponse)
	assert.False(suite.T(), resp.Success)
	assert.Equal(suite.T(), "NO_SLOT_GIVEN", resp.ResultCode)
}

func (suite *HandleEquipSuite) TestNoTarget() {

	suite.msg.GetEquipRequest().SlotLocation = 1

	suite.world.handleEquip(suite.handle)

	assert.Equal(suite.T(), 1, suite.player.SentMessageCount())
	resp := suite.player.GetSentResponse(0).(message.EquipResponse)
	assert.False(suite.T(), resp.Success)
	assert.Equal(suite.T(), "NO_TARGET", resp.ResultCode)
}
