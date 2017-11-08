package world

import (
	"github.com/stretchr/testify/assert"
	"github.com/trasa/watchmud/client"
	"github.com/trasa/watchmud/gameserver"
	"github.com/trasa/watchmud/message"
	"github.com/trasa/watchmud/player"
	"testing"
)

func newWhoRequestHandlerParameter(t *testing.T, c *client.TestClient) *gameserver.HandlerParameter {
	msg, err := message.NewGameMessage(message.WhoRequest{})
	assert.NoError(t, err)
	return gameserver.NewHandlerParameter(c, msg)
}

func TestWho_success(t *testing.T) {
	w := newTestWorld()
	p := player.NewTestPlayer("guy")
	w.AddPlayer(p)
	c := client.NewTestClient(p)

	w.handleWho(newWhoRequestHandlerParameter(t, c))

	assert.Equal(t, 1, p.SentMessageCount())
	resp := p.GetSentResponse(0).(message.WhoResponse)
	assert.True(t, resp.Success)
	assert.Equal(t, 1, len(resp.PlayerInfo))
	assert.Equal(t, "guy", resp.PlayerInfo[0].PlayerName)
	assert.NotEqual(t, "", resp.PlayerInfo[0].ZoneName)
	assert.NotEqual(t, "", resp.PlayerInfo[0].RoomName)
}

func TestWho_notInRoom(t *testing.T) {
	w := newTestWorld()
	p := player.NewTestPlayer("guy")
	w.AddPlayer(p)
	w.playerRooms.Remove(p)
	c := client.NewTestClient(p)

	w.handleWho(newWhoRequestHandlerParameter(t, c))

	assert.Equal(t, 1, p.SentMessageCount())
	resp := p.GetSentResponse(0).(message.WhoResponse)
	assert.True(t, resp.Success)
	assert.Equal(t, "", resp.PlayerInfo[0].ZoneName)
	assert.Equal(t, "", resp.PlayerInfo[0].RoomName)
}

func TestWho_sort(t *testing.T) {
	w := newTestWorld()
	z := player.NewTestPlayer("z")
	y := player.NewTestPlayer("y")
	w.AddPlayer(z, y)
	c := client.NewTestClient(z)

	w.handleWho(newWhoRequestHandlerParameter(t, c))

	assert.Equal(t, "y", z.GetSentResponse(0).(message.WhoResponse).PlayerInfo[0].PlayerName)
	assert.Equal(t, "z", z.GetSentResponse(0).(message.WhoResponse).PlayerInfo[1].PlayerName)
}

func TestWho_logoutRemovesPlayer(t *testing.T) {
	w := newTestWorld()
	z := player.NewTestPlayer("z")
	y := player.NewTestPlayer("y")
	w.AddPlayer(z, y)
	w.RemovePlayer(y)
	c := client.NewTestClient(z)

	w.handleWho(newWhoRequestHandlerParameter(t, c))

	assert.Equal(t, 1, len(z.GetSentResponse(0).(message.WhoResponse).PlayerInfo))
	assert.Equal(t, "z", z.GetSentResponse(0).(message.WhoResponse).PlayerInfo[0].PlayerName)
}
