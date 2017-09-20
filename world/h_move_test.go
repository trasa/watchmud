package world

import (
	"github.com/stretchr/testify/assert"
	"github.com/trasa/watchmud/direction"
	"github.com/trasa/watchmud/message"
	"github.com/trasa/watchmud/player"
	"log"
	"testing"
)

func TestMove_butYouCant(t *testing.T) {
	w := newTestWorld()
	p := player.NewTestPlayer("p")
	w.AddPlayer(p)

	msg := message.IncomingMessage{
		Player: p,
		Request: message.MoveRequest{
			Request:   message.RequestBase{MessageType: "move"},
			Direction: direction.NORTH,
		},
	}

	w.handleMove(&msg)

	log.Printf("%d", p.SentMessageCount())
	if p.SentMessageCount() != 1 {
		t.Fatalf("Expected message count of 1: %d", p.SentMessageCount())
	}
	resp := p.GetSentResponse(0).(message.Response)
	assert.False(t, resp.IsSuccessful())
	assert.Equal(t, resp.GetMessageType(), "move")
	assert.Equal(t, resp.GetResultCode(), "CANT_GO_THAT_WAY")
}
