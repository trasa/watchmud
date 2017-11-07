package world

import (
	"github.com/stretchr/testify/assert"
	"github.com/trasa/watchmud/message"
	"github.com/trasa/watchmud/player"
	"testing"
)

func TestDrop_success(t *testing.T) {
	w := newTestWorld()
	p := player.NewTestPlayer("foo")
	w.AddPlayer(p)

	// get first
	getmsg := message.IncomingMessage{
		Player: p,
		Request: message.GetRequest{
			Request: message.RequestBase{MessageType: "get"},
			Targets: []string{"knife"},
		},
	}
	w.handleGet(&getmsg)

	// now drop
	dropmsg := message.IncomingMessage{
		Player: p,
		Request: message.DropRequest{
			Request: message.RequestBase{MessageType: "drop"},
			Target:  "knife",
		},
	}
	w.handleDrop(&dropmsg)

	assert.Equal(t, 2, p.SentMessageCount())
	getresp := p.GetSentResponse(0).(message.GetResponse)
	assert.True(t, getresp.IsSuccessful())
	dropresp := p.GetSentResponse(1).(message.DropResponse)
	assert.True(t, dropresp.IsSuccessful())

	assert.Equal(t, 0, len(p.GetAllInventory()))
	assert.Equal(t, 1, len(w.startRoom.GetAllInventory()))
}

func TestDrop_NoTarget(t *testing.T) {
	w := newTestWorld()
	p := player.NewTestPlayer("foo")
	w.AddPlayer(p)

	// drop
	dropmsg := message.IncomingMessage{
		Player: p,
		Request: message.DropRequest{
			Request: message.RequestBase{MessageType: "drop"},
			Target:  "",
		},
	}
	w.handleDrop(&dropmsg)

	assert.Equal(t, 1, p.SentMessageCount())
	dropresp := p.GetSentResponse(0).(message.DropResponse)
	assert.False(t, dropresp.IsSuccessful())
	assert.Equal(t, "NO_TARGET", dropresp.GetResultCode())
}

func TestDrop_NotFound(t *testing.T) {
	w := newTestWorld()
	p := player.NewTestPlayer("foo")
	w.AddPlayer(p)

	// drop (but you don't have one)
	dropmsg := message.IncomingMessage{
		Player: p,
		Request: message.DropRequest{
			Request: message.RequestBase{MessageType: "drop"},
			Target:  "knife",
		},
	}
	w.handleDrop(&dropmsg)

	assert.Equal(t, 1, p.SentMessageCount())
	dropresp := p.GetSentResponse(0).(message.DropResponse)
	assert.False(t, dropresp.IsSuccessful())
	assert.Equal(t, "TARGET_NOT_FOUND", dropresp.GetResultCode())
	assert.Equal(t, 0, len(p.GetAllInventory()))
}
