package world

import (
	"github.com/stretchr/testify/assert"
	"github.com/trasa/watchmud/client"
	"github.com/trasa/watchmud/gameserver"
	"github.com/trasa/watchmud/message"
	"github.com/trasa/watchmud/player"
	"testing"
)

func newTellRequestHandlerParameter(t *testing.T, c *client.TestClient, receiver string, value string) *gameserver.HandlerParameter {
	msg, err := message.NewGameMessage(message.TellRequest{
		ReceiverPlayerName: receiver,
		Value:              value,
	})
	assert.NoError(t, err)
	return gameserver.NewHandlerParameter(c, msg)
}

// tell receiver about it
func TestHandleTell_success(t *testing.T) {
	// arrange
	w := newTestWorld()
	senderPlayer := player.NewTestPlayer("sender")
	receiverPlayer := player.NewTestPlayer("receiver")
	w.AddPlayer(receiverPlayer)
	w.AddPlayer(senderPlayer)
	c := client.NewTestClient(senderPlayer)

	// act
	w.handleTell(newTellRequestHandlerParameter(t, c, "receiver", "hi"))

	// assert
	// assert tell to receiver
	assert.Equal(t, 1, receiverPlayer.SentMessageCount())
	recdMessage := receiverPlayer.GetSentResponse(0).(message.TellNotification)
	assert.Equal(t, senderPlayer.GetName(), recdMessage.Sender)
	assert.Equal(t, "hi", recdMessage.Value)

	// assert tell-response to sender
	assert.Equal(t, 1, senderPlayer.SentMessageCount())
	senderResponse := senderPlayer.GetSentResponse(0).(message.TellResponse)
	assert.Equal(t, "OK", senderResponse.ResultCode)
	assert.True(t, senderResponse.Success)
}

func TestHandleTell_receiverNotFound(t *testing.T) {
	// arrange
	w := newTestWorld()
	senderPlayer := player.NewTestPlayer("sender")
	c := client.NewTestClient(senderPlayer)
	// note: receiver doesn't exist
	w.AddPlayer(senderPlayer)

	// act
	w.handleTell(newTellRequestHandlerParameter(t, c, "receiver", "hi"))

	// assert tell-response to sender
	assert.Equal(t, 1, senderPlayer.SentMessageCount())
	senderResponse := senderPlayer.GetSentResponse(0).(message.TellResponse)
	assert.Equal(t, "TO_PLAYER_NOT_FOUND", senderResponse.ResultCode)
	assert.False(t, senderResponse.Success)
}
