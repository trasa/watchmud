package world

import (
	"github.com/trasa/watchmud/message"
	"github.com/trasa/watchmud/player"
	"testing"
)

// tell receiver about it
func TestHandleTell_success(t *testing.T) {
	// arrange
	w := newTestWorld()
	senderPlayer := player.NewTestPlayer("sender")
	receiverPlayer := player.NewTestPlayer("receiver")
	w.AddPlayer(receiverPlayer)
	w.AddPlayer(senderPlayer)

	msg := message.IncomingMessage{
		Player: senderPlayer,
		Request: message.TellRequest{
			Request:            message.RequestBase{MessageType: "tell"},
			ReceiverPlayerName: "receiver",
			Value:              "hi",
		},
	}

	// act
	w.handleTell(&msg)

	// assert
	// assert tell to receiver

	if receiverPlayer.SentMessageCount() != 1 {
		t.Errorf("expected receiver to get a message: %d", receiverPlayer.SentMessageCount())
	}
	recdMessage := receiverPlayer.GetSentResponse(0).(message.TellNotification)
	if recdMessage.Sender != senderPlayer.GetName() {
		t.Errorf("Didn't get expected senderPlayer.Name: %s", recdMessage.Sender)
	}
	if recdMessage.GetMessageType() != "tell_notification" {
		t.Errorf("MsgType wasn't tell: %s", recdMessage.GetMessageType())
	}
	if recdMessage.Value != "hi" {
		t.Errorf("value wasn't hi: %s", recdMessage.Value)
	}

	// assert tell-response to sender
	if senderPlayer.SentMessageCount() != 1 {
		t.Errorf("expected sender to get a response: %d", senderPlayer.SentMessageCount())
	}
	senderResponse := senderPlayer.GetSentResponse(0).(message.Response)
	if senderResponse.GetMessageType() != "tell" {
		t.Errorf("expected sender response tell: %s", senderResponse.GetMessageType())
	}
	if senderResponse.GetResultCode() != "OK" {
		t.Errorf("expected sender response ok: %s", senderResponse.GetResultCode())
	}
	if !senderResponse.IsSuccessful() {
		t.Error("expected sender response to be successful")
	}
}

func TestHandleTell_receiverNotFound(t *testing.T) {
	// arrange
	w := newTestWorld()
	senderPlayer := player.NewTestPlayer("sender")
	// note: receiver doesn't exist
	w.AddPlayer(senderPlayer)

	msg := message.IncomingMessage{
		Player: senderPlayer,
		Request: message.TellRequest{
			Request:            message.RequestBase{MessageType: "tell"},
			ReceiverPlayerName: "receiver",
			Value:              "hi",
		},
	}

	// act
	w.handleTell(&msg)

	// assert tell-response to sender
	if senderPlayer.SentMessageCount() != 1 {
		t.Error("expected sender to get a response")
	}

	senderResponse := senderPlayer.GetSentResponse(0).(message.Response)
	if senderResponse.GetMessageType() != "tell" {
		t.Errorf("sender response message type: %s", senderResponse.GetMessageType())
	}
	if senderResponse.GetResultCode() != "TO_PLAYER_NOT_FOUND" {
		t.Errorf("sender response: %s", senderResponse.GetResultCode())
	}
	if senderResponse.IsSuccessful() {
		t.Error("expected sender response to be a failure")
	}
}
