package world

import (
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
