package world

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/trasa/watchmud/client"
	"github.com/trasa/watchmud/gameserver"
	"github.com/trasa/watchmud/message"
	"github.com/trasa/watchmud/player"
	"testing"
)

type HandleGetSuite struct {
	suite.Suite
	world      *World
	player     *player.TestPlayer
	testClient *client.TestClient
}

func TestHandleGetSuite(t *testing.T) {
	suite.Run(t, new(HandleGetSuite))
}

func (suite *HandleGetSuite) SetupTest() {
	suite.world = newTestWorld()
	suite.player = player.NewTestPlayer("foo")
	suite.world.AddPlayer(suite.player)
	suite.testClient = client.NewTestClient(suite.player)
}

func newGetRequestHandlerParameter(t *testing.T, c *client.TestClient, target string) *gameserver.HandlerParameter {
	msg, err := message.NewGameMessage(message.GetRequest{Target: target})
	assert.NoError(t, err)
	return gameserver.NewHandlerParameter(c, msg)
}

func (suite *HandleGetSuite) TestGet_success() {
	// start off with two items in the room and zero in the player
	suite.Assert().Equal(2, len(suite.world.startRoom.GetAllInventory()))
	suite.Assert().Equal(0, len(suite.player.GetAllInventory()))

	getHP := newGetRequestHandlerParameter(suite.T(), suite.testClient, "knife")
	suite.world.handleGet(getHP)

	assert.Equal(suite.T(), 1, suite.player.SentMessageCount())
	resp := suite.player.GetSentResponse(0).(message.GetResponse)
	assert.True(suite.T(), resp.Success)
	// player has one item
	assert.Equal(suite.T(), 1, len(suite.player.GetAllInventory()))
	foundinv, exists := suite.player.GetInventoryByName("knife")
	assert.True(suite.T(), exists)
	assert.Equal(suite.T(), "knife", foundinv.Definition.Name)

	// there's one other item in the room now
	assert.Equal(suite.T(), 1, len(suite.world.startRoom.GetAllInventory()))
}

func (suite *HandleGetSuite) TestGet_aliasTarget() {
	getHP := newGetRequestHandlerParameter(suite.T(), suite.testClient, "iron")
	suite.world.handleGet(getHP)

	assert.Equal(suite.T(), 1, suite.player.SentMessageCount())
	resp := suite.player.GetSentResponse(0).(message.GetResponse)
	assert.True(suite.T(), resp.Success)
	assert.Equal(suite.T(), 1, len(suite.player.GetAllInventory()))
	foundinv, exists := suite.player.GetInventoryByName("iron_helmet")
	assert.True(suite.T(), exists)
	assert.Equal(suite.T(), "iron_helmet", foundinv.Definition.Name)
	assert.Equal(suite.T(), 1, len(suite.world.startRoom.GetAllInventory()))
}

func (suite *HandleGetSuite) TestGet_targetNotInRoom() {
	getHP := newGetRequestHandlerParameter(suite.T(), suite.testClient, "bag_of_coins")
	suite.world.handleGet(getHP)

	assert.Equal(suite.T(), 1, suite.player.SentMessageCount())
	resp := suite.player.GetSentResponse(0).(message.GetResponse)
	assert.False(suite.T(), resp.Success)
	assert.Equal(suite.T(), "TARGET_NOT_FOUND", resp.GetResultCode())
	// player has zero items still
	assert.Equal(suite.T(), 0, len(suite.player.GetAllInventory()))
	// still two items in start room
	assert.Equal(suite.T(), 2, len(suite.world.startRoom.GetAllInventory()))
}

func (suite *HandleGetSuite) TestGet_noTarget() {
	getHP := newGetRequestHandlerParameter(suite.T(), suite.testClient, "")
	suite.world.handleGet(getHP)

	assert.Equal(suite.T(), 1, suite.player.SentMessageCount())
	resp := suite.player.GetSentResponse(0).(message.GetResponse)
	assert.False(suite.T(), resp.Success)
	assert.Equal(suite.T(), "NO_TARGET", resp.GetResultCode())
	// player has zero items, start room still has 2
	assert.Equal(suite.T(), 0, len(suite.player.GetAllInventory()))
	assert.Equal(suite.T(), 2, len(suite.world.startRoom.GetAllInventory()))
}

func (suite *HandleGetSuite) TestGet_playerAddFail() {
	// TODO: some sort of world-wide list of inventory definitions
	// give the player a knife to start with
	// note that two different objects should not have the same instance id
	// -- this is an arbitrary case to make the test work...
	inv, exists := suite.world.startRoom.GetInventoryByName("knife")
	assert.True(suite.T(), exists)
	suite.player.AddInventory(inv)

	getHP := newGetRequestHandlerParameter(suite.T(), suite.testClient, "knife")
	suite.world.handleGet(getHP)

	assert.Equal(suite.T(), 1, suite.player.SentMessageCount())
	resp := suite.player.GetSentResponse(0).(message.GetResponse)
	assert.False(suite.T(), resp.Success)
	assert.Equal(suite.T(), "ADD_INVENTORY_ERROR", resp.GetResultCode())
	// player just has one item we added at beginning of this method
	assert.Equal(suite.T(), 1, len(suite.player.GetAllInventory()))
	// room still has its two items
	assert.Equal(suite.T(), 2, len(suite.world.startRoom.GetAllInventory()))
}

// TODO: test case for when room.Inventory.Remove fails
// need to figure out how to mock the room
