package world

import (
	"github.com/stretchr/testify/assert"
	"github.com/trasa/watchmud/client"
	"github.com/trasa/watchmud/gameserver"
	"github.com/trasa/watchmud/message"
	"github.com/trasa/watchmud/player"
	"testing"
)

func TestWorld_handleMessage_unknownMessageType(t *testing.T) {
	w := newTestWorld()
	p := player.NewTestPlayer("sender")
	c := client.NewTestClient(p)
	m := &message.GameMessage{} // not a valid message, has no inner type
	h := gameserver.NewHandlerParameter(c, m)

	w.HandleIncomingMessage(h)

	resp := p.GetSentResponse(0).(message.ErrorResponse)
	assert.False(t, resp.Success)
	assert.Equal(t, "UNKNOWN_MESSAGE_TYPE", resp.ResultCode)
}

func TestWorld_RemovePlayer(t *testing.T) {
	w := newTestWorld()
	p := player.NewTestPlayer("dood")
	w.AddPlayer(p)
	w.RemovePlayer(p)

	assert.Equal(t, 0, w.playerList.Count())
	assert.Nil(t, w.playerRooms.playerToRoom[p])
	assert.Equal(t, 0, len(w.playerRooms.roomToPlayers.Get(w.startRoom)))
	assert.Equal(t, 0, len(w.startRoom.GetPlayers()))
}
