package world

import (
	"github.com/stretchr/testify/assert"
	"github.com/trasa/watchmud/client"
	"github.com/trasa/watchmud/direction"
	"github.com/trasa/watchmud/gameserver"
	"github.com/trasa/watchmud/message"
	"github.com/trasa/watchmud/player"
	"log"
	"testing"
)

func newMoveRequestHandlerParameter(t *testing.T, c *client.TestClient, dir direction.Direction) *gameserver.HandlerParameter {
	msg, err := message.NewGameMessage(message.MoveRequest{Direction: int32(dir)})
	assert.NoError(t, err)
	return gameserver.NewHandlerParameter(c, msg)
}

func TestMove_butYouCant(t *testing.T) {
	w := newTestWorld()
	p := player.NewTestPlayer("p")
	w.AddPlayer(p)
	c := client.NewTestClient(p)

	w.handleMove(newMoveRequestHandlerParameter(t, c, direction.NORTH))

	log.Printf("%d", p.SentMessageCount())
	if p.SentMessageCount() != 1 {
		t.Fatalf("Expected message count of 1: %d", p.SentMessageCount())
	}
	resp := p.GetSentResponse(0).(message.MoveResponse)
	assert.False(t, resp.Success)
	assert.Equal(t, resp.ResultCode, "CANT_GO_THAT_WAY")
}
