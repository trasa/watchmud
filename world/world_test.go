package world

import (
	"github.com/trasa/watchmud/message"
	"github.com/trasa/watchmud/player"
	"testing"
)

// create a new test world
func newTestWorld() *World {

	testZone := Zone{
		Id:    "sample",
		Name:  "Sample Zone",
		Rooms: make(map[string]*Room),
	}
	w := &World{
		zones:       make(map[string]*Zone),
		playerList:  player.NewList(),
		playerRooms: NewPlayerRoomMap(),
	}
	w.zones[testZone.Id] = &testZone

	testRoom := NewRoom(&testZone, "start", "Test Room", "this is a test room.")
	testZone.Rooms[testRoom.Id] = testRoom
	w.startRoom = testRoom
	return w
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
