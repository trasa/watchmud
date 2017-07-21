package world

import (
	"github.com/trasa/watchmud/message"
	"github.com/trasa/watchmud/player"
	"testing"
)

// TODO this breaks because of channels!

// create a new test world
func newTestWorld() *World {
	return &World{
		PlayerList: player.NewList(),
	}
}

func TestWorld_handleMessage_unknownMessageType(t *testing.T) {
	w := newTestWorld()
	p := NewTestPlayer("sender")

	msg := message.IncomingMessage{
		Player: p,
		Request: &message.RequestBase{
			MessageType: "asdfasdf",
		},
	}
	w.HandleIncomingMessage(&msg)

	resp := p.GetSentResponse(0)
	if resp.Successful {
		t.Error("should not succeed")
	}
	if resp.ResultCode != "UNKNOWN_MESSAGE_TYPE" {
		t.Errorf("unexpected result code: %s", resp.ResultCode)
	}
}
