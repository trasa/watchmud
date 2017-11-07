package world

import (
	"github.com/stretchr/testify/assert"
	"github.com/trasa/watchmud/message"
	"github.com/trasa/watchmud/player"
	"testing"
)

func TestWorld_handleMessage_unknownMessageType(t *testing.T) {
	w := newTestWorld()
	p := player.NewTestPlayer("sender")

	msg := message.IncomingMessage{
		Player: p,
		Request: &message.RequestBase{
			MessageType: "asdfasdf",
		},
	}
	w.HandleIncomingMessage(&msg)

	resp := p.GetSentResponse(0).(message.Response)
	if resp.IsSuccessful() {
		t.Error("should not succeed")
	}
	if resp.GetResultCode() != "UNKNOWN_MESSAGE_TYPE" {
		t.Errorf("unexpected result code: %s", resp.GetResultCode())
	}
}

func TestWorld_RemovePlayer(t *testing.T) {
	w := newTestWorld()
	p := player.NewTestPlayer("dood")
	w.AddPlayer(p)
	w.RemovePlayer(p)

	assert.Equal(t, 0, w.playerList.Count())
	assert.Nil(t, w.playerRooms.playerToRoom[p])
	assert.Equal(t, 0, len(w.playerRooms.roomToPlayers.Get(w.startRoom)))
	assert.Equal(t, 0, len(w.startRoom.GetPlayers()))
}
