package world

import (
	"github.com/stretchr/testify/assert"
	"github.com/trasa/watchmud/client"
	"github.com/trasa/watchmud/gameserver"
	"github.com/trasa/watchmud/message"
	"github.com/trasa/watchmud/player"
	"testing"
)

func newGetRequestHandlerParameter(t *testing.T, c *client.TestClient, target string) *gameserver.HandlerParameter {
	return newGetRequestHandlerParameterWithSlice(t, c, []string{target})
}

func newGetRequestHandlerParameterWithSlice(t *testing.T, c *client.TestClient, targets []string) *gameserver.HandlerParameter {
	msg, err := message.NewGameMessage(message.GetRequest{Targets: targets})
	assert.NoError(t, err)
	return gameserver.NewHandlerParameter(c, msg)
}

func TestGet_success(t *testing.T) {
	w := newTestWorld()
	p := player.NewTestPlayer("foo")
	testClient := client.NewTestClient(p)
	w.AddPlayer(p)

	getHP := newGetRequestHandlerParameter(t, testClient, "knife")
	w.handleGet(getHP)

	assert.Equal(t, 1, p.SentMessageCount())
	resp := p.GetSentResponse(0).(message.GetResponse)
	assert.True(t, resp.Success)
	assert.Equal(t, 1, len(p.GetAllInventory()))
	foundinv, exists := p.GetInventoryByName("knife")
	assert.True(t, exists)
	assert.Equal(t, "knife", foundinv.Definition.Name)
	assert.Equal(t, 0, len(w.startRoom.GetAllInventory()))
}

func TestGet_targetNotInRoom(t *testing.T) {
	w := newTestWorld()
	p := player.NewTestPlayer("foo")
	testClient := client.NewTestClient(p)
	w.AddPlayer(p)

	getHP := newGetRequestHandlerParameter(t, testClient, "bag_of_coins")
	w.handleGet(getHP)

	assert.Equal(t, 1, p.SentMessageCount())
	resp := p.GetSentResponse(0).(message.GetResponse)
	assert.False(t, resp.Success)
	assert.Equal(t, "TARGET_NOT_FOUND", resp.GetResultCode())
	assert.Equal(t, 0, len(p.GetAllInventory()))
	assert.Equal(t, 1, len(w.startRoom.GetAllInventory()))
}

func TestGet_noTarget(t *testing.T) {
	w := newTestWorld()
	p := player.NewTestPlayer("foo")
	testClient := client.NewTestClient(p)
	w.AddPlayer(p)

	getHP := newGetRequestHandlerParameterWithSlice(t, testClient, []string{})
	w.handleGet(getHP)

	assert.Equal(t, 1, p.SentMessageCount())
	resp := p.GetSentResponse(0).(message.GetResponse)
	assert.False(t, resp.Success)
	assert.Equal(t, "NO_TARGET", resp.GetResultCode())
	assert.Equal(t, 0, len(p.GetAllInventory()))
	assert.Equal(t, 1, len(w.startRoom.GetAllInventory()))
}

func TestGet_playerAddFail(t *testing.T) {
	w := newTestWorld()
	p := player.NewTestPlayer("foo")
	c := client.NewTestClient(p)
	w.AddPlayer(p)

	// TODO: some sort of world-wide list of inventory definitions
	// give the player a knife to start with
	// note that two different objects should not have the same instance id
	// -- this is an arbitrary case to make the test work...
	inv, exists := w.startRoom.GetInventoryByName("knife")
	assert.True(t, exists)
	p.AddInventory(inv)

	getHP := newGetRequestHandlerParameter(t, c, "knife")
	w.handleGet(getHP)

	assert.Equal(t, 1, p.SentMessageCount())
	resp := p.GetSentResponse(0).(message.GetResponse)
	assert.False(t, resp.Success)
	assert.Equal(t, "ADD_INVENTORY_ERROR", resp.GetResultCode())
	assert.Equal(t, 1, len(p.GetAllInventory()))
	assert.Equal(t, 1, len(w.startRoom.GetAllInventory()))
}

// TODO: test case for when room.Inventory.Remove fails
// need to figure out how to mock the room
