package world

import (
	"github.com/trasa/watchmud/message"
	"github.com/trasa/watchmud/response"
	"testing"
)

// tell receiver about it
func TestHandleTell_success(t *testing.T) {
	// arrange
	w := newTestWorld()
	senderPlayer := NewTestPlayer("sender")
	receiverPlayer := NewTestPlayer("receiver")
	w.AddPlayer(receiverPlayer)
	w.AddPlayer(senderPlayer)

	msg := message.IncomingMessage{
		Player: senderPlayer,
		Body:   make(map[string]string),
	}
	msg.Body["to"] = "receiver"
	msg.Body["value"] = "hi"

	// act
	w.handleTell(&msg)

	// assert
	// assert tell to receiver

	if len(receiverPlayer.client.(*TestClient).tosend) != 1 {
		t.Errorf("expected receiver to get a message %s", receiverPlayer.client.(*TestClient).tosend)
	}
	recdMessage := receiverPlayer.client.(*TestClient).tosend[0].(response.TellNotification)
	if recdMessage.From != senderPlayer.name {
		t.Errorf("Didn't get expected senderPlayer.Name: %s", recdMessage.From)
	}
	if recdMessage.MessageType != "tell" {
		t.Errorf("MsgType wasn't tell: %s", recdMessage.MessageType)
	}
	if recdMessage.Value != "hi" {
		t.Errorf("value wasn't hi: %s", recdMessage.Value)
	}

	// assert tell-response to sender
	if len(senderPlayer.client.(*TestClient).tosend) != 1 {
		t.Errorf("expected sender to get a response %s", senderPlayer.client.(*TestClient).tosend)
	}
	senderResponse := senderPlayer.client.(*TestClient).tosend[0].(response.Response)
	if senderResponse.MessageType != "tell" {
		t.Errorf("expected sender response tell: %s", senderResponse.MessageType)
	}
	if senderResponse.ResultCode != "OK" {
		t.Errorf("expected sender response ok: %s", senderResponse.ResultCode)
	}
	if !senderResponse.Successful {
		t.Error("expected sender response to be successful")
	}
}

func TestHandleTell_receiverNotFound(t *testing.T) {
	// arrange
	w := newTestWorld()
	senderPlayer := NewTestPlayer("sender")
	// note: receiver doesn't exist
	w.AddPlayer(senderPlayer)

	msg := message.IncomingMessage{
		Player: senderPlayer,
		Body:   make(map[string]string),
	}
	msg.Body["to"] = "receiver"
	msg.Body["value"] = "hi"

	// act
	w.handleTell(&msg)

	// assert tell-response to sender
	if len(senderPlayer.client.(*TestClient).tosend) != 1 {
		t.Error("expected sender to get a response")
	}

	senderResponse := senderPlayer.client.(*TestClient).tosend[0].(response.Response)
	if senderResponse.MessageType != "tell" {
		t.Errorf("sender response message type: %s", senderResponse.MessageType)
	}
	if senderResponse.ResultCode != "TO_PLAYER_NOT_FOUND" {
		t.Errorf("sender response: %s", senderResponse.ResultCode)
	}
	if senderResponse.Successful {
		t.Error("expected sender response to be a failure")
	}
}
