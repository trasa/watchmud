package world

import (
	"github.com/stretchr/testify/assert"
	"github.com/trasa/watchmud/direction"
	"github.com/trasa/watchmud/message"
	"testing"
)

func TestGo_butYouCant(t *testing.T) {
	w := newTestWorld()
	p := NewTestPlayer("p")
	w.AddPlayer(p)

	msg := message.IncomingMessage{
		Player: p,
		Request: message.GoRequest{
			Request:   message.RequestBase{MessageType: "go"},
			Direction: direction.NORTH,
		},
	}

	w.handleGo(&msg)

	if len(p.sent) != 1 {
		t.Errorf("Expected message %s", p.sent)
	}
	resp := p.sent[0].(message.Response)
	assert.False(t, resp.Successful)
	assert.Equal(t, resp.MessageType, "go")
	assert.Equal(t, resp.ResultCode, "CANT_GO_THAT_WAY")
}
