package world

import (
	"github.com/stretchr/testify/assert"
	"github.com/trasa/watchmud-message"
	"github.com/trasa/watchmud/client"
	"github.com/trasa/watchmud/gameserver"
	"github.com/trasa/watchmud/player"
	"testing"
)

func newTellAllRequestHandlerParameter(t *testing.T, c *client.TestClient, value string) *gameserver.HandlerParameter {
	msg, err := message.NewGameMessage(message.TellAllRequest{
		Value: value,
	})
	assert.NoError(t, err)
	return gameserver.NewHandlerParameter(c, msg)
}

func TestHandleTellAll_success(t *testing.T) {
	w := newTestWorld()
	senderPlayer := player.NewTestPlayer("sender")
	otherPlayer := player.NewTestPlayer("other")
	bobPlayer := player.NewTestPlayer("bob")
	w.AddPlayer(senderPlayer, otherPlayer, bobPlayer)
	c := client.NewTestClient(senderPlayer)

	w.handleTellAll(newTellAllRequestHandlerParameter(t, c, "hi"))

	// did we tell otherPlayer?
	assert.Equal(t, 1, otherPlayer.SentMessageCount())
	assert.Equal(t, 1, bobPlayer.SentMessageCount())
	// sender should have gotten response but NOT part of the send to all players
	assert.Equal(t, 1, senderPlayer.SentMessageCount())
	senderResponse := senderPlayer.GetSentResponse(0).(message.TellAllResponse)
	assert.True(t, senderResponse.Success)
}

func TestHandleTellAll_noValue(t *testing.T) {
	w := newTestWorld()
	senderPlayer := player.NewTestPlayer("sender")
	otherPlayer := player.NewTestPlayer("other")
	w.AddPlayer(senderPlayer, otherPlayer)
	c := client.NewTestClient(senderPlayer)

	w.handleTellAll(newTellAllRequestHandlerParameter(t, c, ""))

	// did we tell otherPlayer? (should be 0)
	assert.Equal(t, 0, otherPlayer.SentMessageCount())

	// sender should have gotten response but NOT part of the send to all players
	assert.Equal(t, 1, senderPlayer.SentMessageCount())

	senderResponse := senderPlayer.GetSentResponse(0).(message.TellAllResponse)
	assert.False(t, senderResponse.Success)
}
