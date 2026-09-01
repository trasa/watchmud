package world

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/trasa/watchmud-message"
	"github.com/trasa/watchmud/client"
	"github.com/trasa/watchmud/gameserver"
	"github.com/trasa/watchmud/loader"
	"github.com/trasa/watchmud/player"
	"github.com/trasa/watchmud/rules"
	"github.com/trasa/watchmud/spaces"
	"github.com/trasa/watchmud/zonereset"
)

func TestWorld_New(t *testing.T) {
	settings := loader.Settings{
		StartZone: "z",
		StartRoom: "r",
		VoidZone:  "z",
		VoidRoom:  "r",
	}
	var species []*rules.Species
	var classes []*rules.Class
	catalog, err := rules.NewCatalog(species, classes)
	assert.NoError(t, err)
	content := loader.NewContent(&settings, catalog)
	content.Zones["z"] = spaces.NewZone("z", "Z", zonereset.NEVER, time.Duration(0))
	content.Zones["z"].AddRoom(spaces.NewRoom(content.Zones["z"], "r", "R", "Description"))

	w, err := New(content)
	assert.NoError(t, err)
	assert.NotNil(t, w)
}

func TestWorld_handleMessage_unknownMessageType(t *testing.T) {
	w, _ := newTestWorld()
	p := player.NewTestPlayer("sender")
	c := client.NewTestClient(p)
	m := &message.GameMessage{} // not a valid message, has no inner type
	h := gameserver.NewHandlerParameter(c, m)

	w.HandleIncomingMessage(h)

	resp := c.GetSentResponse(0).(message.ErrorResponse)
	assert.False(t, resp.Success)
	assert.Equal(t, "UNKNOWN_MESSAGE_TYPE", resp.ResultCode)
}

func TestWorld_RemovePlayer(t *testing.T) {
	w, _ := newTestWorld()
	p := player.NewTestPlayer("dood")
	w.AddPlayer(p)
	w.RemovePlayer(p)

	assert.Equal(t, 0, w.playerList.Count())
	assert.Nil(t, w.playerRooms.playerToRoom[p])
	assert.Equal(t, 0, len(w.playerRooms.roomToPlayers.Get(w.StartRoom)))
	assert.Equal(t, 0, len(w.StartRoom.GetPlayers()))
}
