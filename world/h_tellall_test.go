package world

import (
	"github.com/trasa/watchmud/message"
	"github.com/trasa/watchmud/player"
	"testing"
)

func TestHandleTellAll_success(t *testing.T) {
	w := newTestWorld()
	senderPlayer := player.NewTestPlayer("sender")
	otherPlayer := player.NewTestPlayer("other")
	bobPlayer := player.NewTestPlayer("bob")
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
	if otherPlayer.SentMessageCount() != 1 {
		t.Error("otherPlayer should have gotten a message")
	}
	if bobPlayer.SentMessageCount() != 1 {
		t.Error("bob should have gotten a message")
	}
	// sender should have gotten response but NOT part of the send to all players
	if senderPlayer.SentMessageCount() != 1 {
		t.Fatalf("sender received wrong number of messages: %d", senderPlayer.SentMessageCount())
	}
	senderResponse := senderPlayer.GetSentResponse(0).(message.Response)
	if senderResponse.GetMessageType() != "tell_all" {
		t.Errorf("incorrect sender response message type: %s", senderResponse.GetMessageType())
	}
	if !senderResponse.IsSuccessful() {
		t.Error("sender response is not successful")
	}
}

func TestHandleTellAll_noValue(t *testing.T) {
	w := newTestWorld()
	senderPlayer := player.NewTestPlayer("sender")
	otherPlayer := player.NewTestPlayer("other")
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
	if otherPlayer.SentMessageCount() != 0 {
		t.Error("no send")
	}
	// sender should have gotten response but NOT part of the send to all players
	if senderPlayer.SentMessageCount() != 1 {
		t.Fatalf("sender received wrong number of messages: %d", senderPlayer.SentMessageCount())
	}
	senderResponse := senderPlayer.GetSentResponse(0).(message.Response)
	if senderResponse.GetMessageType() != "tell_all" {
		t.Errorf("incorrect sender response message type: %s", senderResponse.GetMessageType())
	}
	if senderResponse.IsSuccessful() {
		t.Error("sender response should fail")
	}
}
