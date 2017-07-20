package server

import (
	"github.com/trasa/watchmud/client"
	"github.com/trasa/watchmud/message"
	"github.com/trasa/watchmud/player"
	"log"
	"testing"
)

// see https://play.golang.org/p/zPLyr3ZOM0 (first attempt)
// then see https://play.golang.org/p/z5athD5fV3 (client is an interface, but now pointer woes)

// create a new test world
func NewTestWorld() *World {
	return &World{
		knownPlayersByName: make(map[string]player.Player),
	}
}

type TestPlayer struct {
	name   string
	client client.Client
}

func (this *TestPlayer) Send(message interface{}) {
	log.Printf("sending for player %s", message)
	this.client.Send(message)
}

func (this *TestPlayer) GetName() string {
	return this.name
}

// create a new test player that can track sent messages through 'sentmessages'
func NewTestPlayer(name string) *TestPlayer {
	c := &TestClient{}
	p := &TestPlayer{
		name:   name,
		client: c,
	}
	c.Player = p
	return p
}

// tell receiver about it
func TestHandleTell_success(t *testing.T) {
	// arrange
	w := NewTestWorld()
	senderPlayer := NewTestPlayer("sender")
	receiverPlayer := NewTestPlayer("receiver")
	w.knownPlayersByName[receiverPlayer.name] = receiverPlayer
	w.knownPlayersByName[senderPlayer.name] = senderPlayer

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
	recdMessage := receiverPlayer.client.(*TestClient).tosend[0].(TellNotification)
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
	senderResponse := senderPlayer.client.(*TestClient).tosend[0].(Response)
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
	w := NewTestWorld()
	senderPlayer := NewTestPlayer("sender")
	// note: receiver doesn't exist
	w.knownPlayersByName[senderPlayer.name] = senderPlayer

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

	senderResponse := senderPlayer.client.(*TestClient).tosend[0].(Response)
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
