package world

import (
	"github.com/stretchr/testify/assert"
	"github.com/trasa/watchmud/message"
	"github.com/trasa/watchmud/player"
	"testing"
)

func TestLook_successful(t *testing.T) {
	w := newTestWorld()
	p := player.NewTestPlayer("testdood")
	other := player.NewTestPlayer("other")
	w.AddPlayer(p)
	w.AddPlayer(other)

	msg := message.IncomingMessage{
		Player: p,
		Request: message.LookRequest{
			Request: message.RequestBase{MessageType: "look"},
		},
	}

	w.HandleIncomingMessage(&msg)

	resp := p.GetSentResponse(0).(message.LookResponse)
	assert.Equal(t, "look", resp.GetMessageType(), "wrong message type")
	assert.True(t, resp.IsSuccessful())
	assert.NotNil(t, resp.RoomDescription.Name)
	assert.NotNil(t, resp.RoomDescription.Description)
	assert.Equal(t, 1, len(resp.RoomDescription.Players))
	assert.Equal(t, "other", resp.RoomDescription.Players[0])
}
