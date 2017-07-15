package server

import (
	"log"
	"testing"
)

// keeps track of a messages sent to players
var sentMessages = make(map[*Player][]interface{})

// create a new test world
func NewTestWorld() *World {
	return &World{
		knownPlayersByName: make(map[string]*Player),
	}
}

// create a new test player that can track sent messages through 'sentmessages'
func NewTestPlayer(name string) *Player {
	p := Player{
		Name: name,
	}
	p.Send = func(message interface{}) {
		log.Printf("sending fake! %s p is %s", message, p.Name)
		sentMessages[&p] = append(sentMessages[&p], message)
	}
	return &p
}

// tell receiver about it
func TestHandleTell_success(t *testing.T) {
	// arrange
	w := NewTestWorld()
	senderPlayer := NewTestPlayer("sender")
	receiverPlayer := NewTestPlayer("receiver")
	w.knownPlayersByName[receiverPlayer.Name] = receiverPlayer
	w.knownPlayersByName[senderPlayer.Name] = senderPlayer

	msg := IncomingMessage{
		Player: senderPlayer,
		Body:   make(map[string]string),
	}
	msg.Body["to"] = "receiver"
	msg.Body["value"] = "hi"

	// act
	w.handleTell(&msg)

	// assert tell to receiver
	if len(sentMessages[receiverPlayer]) != 1 {
		t.Errorf("expected receiver to get a message %s", sentMessages)
	}
	recdMessage := sentMessages[receiverPlayer][0].(TellNotification)
	if recdMessage.From != senderPlayer.Name {
		t.Errorf("Didn't get expected senderPlayer.Name: %s", recdMessage.From)
	}
	if recdMessage.MessageType != "tell" {
		t.Errorf("MsgType wasn't tell: %s", recdMessage.MessageType)
	}
	if recdMessage.Value != "hi" {
		t.Errorf("value wasn't hi: %s", recdMessage.Value)
	}

	// assert tell-response to sender
	if len(sentMessages[senderPlayer]) != 1 {
		t.Errorf("expected sender to get a response %s", sentMessages)
	}
	senderResponse := sentMessages[senderPlayer][0].(Response)
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
