package world

import (
	"github.com/stretchr/testify/assert"
	"github.com/trasa/watchmud/message"
	"testing"
)

func TestWho_success(t *testing.T) {
	w := newTestWorld()
	p := NewTestPlayer("guy")
	w.AddPlayer(p)
	msg := message.IncomingMessage{
		Player: p,
		Request: message.WhoRequest{
			Request: message.RequestBase{MessageType: "who"},
		},
	}

	w.handleWho(&msg)

	assert.Equal(t, 1, len(p.sent))
	resp := p.sent[0].(message.WhoResponse)
	assert.True(t, resp.Successful)
	assert.Equal(t, 1, len(resp.PlayerInfo))
	assert.Equal(t, "guy", resp.PlayerInfo[0].PlayerName)
	assert.NotEqual(t, "", resp.PlayerInfo[0].ZoneName)
	assert.NotEqual(t, "", resp.PlayerInfo[0].RoomName)
}

func TestWho_notInRoom(t *testing.T) {
	w := newTestWorld()
	p := NewTestPlayer("guy")
	w.AddPlayer(p)
	w.playerRooms.Remove(p)

	msg := message.IncomingMessage{
		Player: p,
		Request: message.WhoRequest{
			Request: message.RequestBase{MessageType: "who"},
		},
	}

	w.handleWho(&msg)

	assert.Equal(t, 1, len(p.sent))
	resp := p.sent[0].(message.WhoResponse)
	assert.True(t, resp.Successful)
	assert.Equal(t, "", resp.PlayerInfo[0].ZoneName)
	assert.Equal(t, "", resp.PlayerInfo[0].RoomName)
}

func TestWho_sort(t *testing.T) {
	w := newTestWorld()
	z := NewTestPlayer("z")
	y := NewTestPlayer("y")
	w.AddPlayer(z, y)

	msg := message.IncomingMessage{
		Player: z,
		Request: message.WhoRequest{
			Request: message.RequestBase{MessageType: "who"},
		},
	}

	w.handleWho(&msg)

	assert.Equal(t, "y", z.sent[0].(message.WhoResponse).PlayerInfo[0].PlayerName)
	assert.Equal(t, "z", z.sent[0].(message.WhoResponse).PlayerInfo[1].PlayerName)

}
