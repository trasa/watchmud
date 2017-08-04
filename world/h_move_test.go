package world

import (
	"github.com/stretchr/testify/assert"
	"github.com/trasa/watchmud/direction"
	"github.com/trasa/watchmud/message"
	"log"
	"testing"
)

func TestMove_butYouCant(t *testing.T) {
	w := newTestWorld()
	p := NewTestPlayer("p")
	w.AddPlayer(p)

	msg := message.IncomingMessage{
		Player: p,
		Request: message.MoveRequest{
			Request:   message.RequestBase{MessageType: "move"},
			Direction: direction.NORTH,
		},
	}

	w.handleMove(&msg)

	log.Printf("%d", len(p.sent))
	if len(p.sent) != 1 {
		t.Fatalf("Expected message %s", p.sent)
	}
	resp := p.sent[0].(message.Response)
	assert.False(t, resp.Successful)
	assert.Equal(t, resp.MessageType, "move")
	assert.Equal(t, resp.ResultCode, "CANT_GO_THAT_WAY")
}
