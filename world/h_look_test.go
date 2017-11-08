package world

import (
	"github.com/stretchr/testify/assert"
	"github.com/trasa/watchmud/client"
	"github.com/trasa/watchmud/gameserver"
	"github.com/trasa/watchmud/message"
	"github.com/trasa/watchmud/player"
	"testing"
)

func newLookRequestHandlerParameter(t *testing.T, c *client.TestClient) *gameserver.HandlerParameter {
	msg, err := message.NewGameMessage(message.LookRequest{})
	assert.NoError(t, err)
	return gameserver.NewHandlerParameter(c, msg)
}

func TestLook_successful(t *testing.T) {
	w := newTestWorld()
	p := player.NewTestPlayer("testdood")
	other := player.NewTestPlayer("other")
	w.AddPlayer(p)
	w.AddPlayer(other)
	c := client.NewTestClient(p)

	w.HandleIncomingMessage(newLookRequestHandlerParameter(t, c))

	resp := p.GetSentResponse(0).(message.LookResponse)
	assert.True(t, resp.Success)
	assert.NotNil(t, resp.RoomDescription.Name)
	assert.NotNil(t, resp.RoomDescription.Description)
	assert.Equal(t, 1, len(resp.RoomDescription.Players))
	assert.Equal(t, "other", resp.RoomDescription.Players[0])
}
