package world

import (
	"github.com/trasa/watchmud/message"
	"testing"
)

func TestHandleTellAll_success(t *testing.T) {
	w := newTestWorld()
	senderPlayer := NewTestPlayer("sender")
	otherPlayer := NewTestPlayer("other")
	bobPlayer := NewTestPlayer("bob")
	w.AddPlayer(senderPlayer, otherPlayer, bobPlayer)

	msg := message.IncomingMessage{
		Player: senderPlayer,
		Request: message.TellAllRequest{
			Request: message.RequestBase{MessageType: "tell_all"},
			Value:   "hi",
		},
	}

	w.handleTellAll(&msg)

	// did we tell otherPlayer?
	if len(otherPlayer.sent) != 1 {
		t.Error("otherPlayer should have gotten a message")
	}
	if len(bobPlayer.sent) != 1 {
		t.Error("bob should have gotten a message")
	}
	// sender should have gotten response but NOT part of the send to all players
	if len(senderPlayer.sent) != 1 {
		t.Fatalf("sender received wrong number of messages: %d", len(senderPlayer.sent))
	}
	senderResponse := senderPlayer.sent[0].(message.Response)
	if senderResponse.MessageType != "tell_all" {
		t.Errorf("incorrect sender response message type: %s", senderResponse.MessageType)
	}
	if !senderResponse.Successful {
		t.Error("sender response is not successful")
	}
}

func TestHandleTellAll_noValue(t *testing.T) {
	w := newTestWorld()
	senderPlayer := NewTestPlayer("sender")
	otherPlayer := NewTestPlayer("other")
	w.AddPlayer(senderPlayer, otherPlayer)

	msg := message.IncomingMessage{
		Player: senderPlayer,
		Request: message.TellAllRequest{
			Request: message.RequestBase{MessageType: "tell_all"},
			Value:   "",
		},
	}

	w.handleTellAll(&msg)

	// did we tell otherPlayer?
	if len(otherPlayer.sent) != 0 {
		t.Error("no send")
	}
	// sender should have gotten response but NOT part of the send to all players
	if len(senderPlayer.sent) != 1 {
		t.Fatalf("sender received wrong number of messages: %d", len(senderPlayer.sent))
	}
	senderResponse := senderPlayer.sent[0].(message.Response)
	if senderResponse.MessageType != "tell_all" {
		t.Errorf("incorrect sender response message type: %s", senderResponse.MessageType)
	}
	if senderResponse.Successful {
		t.Error("sender response should fail")
	}
}
