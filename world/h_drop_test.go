package world

import (
	"github.com/stretchr/testify/assert"
	"github.com/trasa/watchmud/client"
	"github.com/trasa/watchmud/gameserver"
	"github.com/trasa/watchmud/message"
	"github.com/trasa/watchmud/player"
	"testing"
)

func TestDrop_success(t *testing.T) {
	w := newTestWorld()
	p := player.NewTestPlayer("foo")
	w.AddPlayer(p)
	testClient := client.NewTestClient(p)

	// get first
	getGameMessage, err := message.NewGameMessage(
		message.GetRequest{
			Targets: []string{"knife"},
		})
	assert.NoError(t, err)
	getHP := gameserver.NewHandlerParameter(testClient, getGameMessage)
	w.handleGet(getHP)

	// now drop
	dropGameMessage, err := message.NewGameMessage(
		message.DropRequest{
			Target: "knife",
		},
	)
	assert.NoError(t, err)
	dropHP := gameserver.NewHandlerParameter(testClient, dropGameMessage)
	w.handleDrop(dropHP)

	assert.Equal(t, 2, p.SentMessageCount())
	getresp := p.GetSentResponse(0).(message.GetResponse)
	assert.True(t, getresp.Success)
	dropresp := p.GetSentResponse(1).(message.DropResponse)
	assert.True(t, dropresp.Success)

	assert.Equal(t, 0, len(p.GetAllInventory()))
	assert.Equal(t, 1, len(w.startRoom.GetAllInventory()))
}

func TestDrop_NoTarget(t *testing.T) {
	w := newTestWorld()
	p := player.NewTestPlayer("foo")
	w.AddPlayer(p)
	testClient := client.NewTestClient(p)

	// drop
	dropGameMessage, err := message.NewGameMessage(message.DropRequest{Target: ""})
	assert.NoError(t, err)
	dropHP := gameserver.NewHandlerParameter(testClient, dropGameMessage)
	w.handleDrop(dropHP)

	assert.Equal(t, 1, p.SentMessageCount())
	dropresp := p.GetSentResponse(0).(message.DropResponse)
	assert.False(t, dropresp.Success)
	assert.Equal(t, "NO_TARGET", dropresp.GetResultCode())
}

func TestDrop_NotFound(t *testing.T) {
	w := newTestWorld()
	p := player.NewTestPlayer("foo")
	w.AddPlayer(p)
	testClient := client.NewTestClient(p)

	// drop (but you don't have one)
	dropGameMessage, err := message.NewGameMessage(message.DropRequest{Target: "knife"})
	assert.NoError(t, err)
	dropHP := gameserver.NewHandlerParameter(testClient, dropGameMessage)
	w.handleDrop(dropHP)

	assert.Equal(t, 1, p.SentMessageCount())
	dropresp := p.GetSentResponse(0).(message.DropResponse)
	assert.False(t, dropresp.Success)
	assert.Equal(t, "TARGET_NOT_FOUND", dropresp.GetResultCode())
	assert.Equal(t, 0, len(p.GetAllInventory()))
}
