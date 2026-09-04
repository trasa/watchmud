package world

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/trasa/watchmud-message"
	"github.com/trasa/watchmud/client"
	"github.com/trasa/watchmud/gameserver"
	"github.com/trasa/watchmud/loader"
	"github.com/trasa/watchmud/player"
	"github.com/trasa/watchmud/rules"
	"github.com/trasa/watchmud/spaces"
	"github.com/trasa/watchmud/zonereset"
)

type worldTestSuite struct {
	suite.Suite
	w *World
	r *player.Recorder
	p *player.Player
	c *client.TestClient
}

// sent returns the ith message in the recorder
func sent[T any](t *testing.T, r *player.Recorder, i int) T {
	t.Helper()
	require.Greater(t, len(r.Sent), i, "wanted Sent[%d], only %d sent", i, len(r.Sent))
	v, ok := r.Sent[i].(T)
	require.Truef(t, ok, "Sent[%d] is %T, want %T", i, r.Sent[i], *new(T))
	return v
}

func (s *worldTestSuite) handlerParameter(req interface{}) *gameserver.HandlerParameter {
	s.T().Helper()
	msg, err := message.NewGameMessage(req)
	s.Require().NoError(err)
	return gameserver.NewHandlerParameter(s.c, msg)
}

func (s *worldTestSuite) SetupTest() {
	s.w, _ = newTestWorld()
	s.r = &player.Recorder{}
	s.p = player.NewTestPlayer("testdood", "testdood", s.r)
	s.w.AddPlayer(s.p)
	s.c = client.NewTestClient(s.p)
}

func TestWorldNew(t *testing.T) {
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

func (s *worldTestSuite) TestUnknownMessageType() {
	m := &message.GameMessage{} // not a valid message, has no inner type
	h := gameserver.NewHandlerParameter(s.c, m)

	s.w.HandleIncomingMessage(h)

	resp := sent[message.ErrorResponse](s.T(), s.r, 0)
	s.Assert().False(resp.Success)
	s.Assert().Equal("UNKNOWN_MESSAGE_TYPE", resp.ResultCode)
}

func (s *worldTestSuite) TestRemovePlayer() {
	s.w.RemovePlayer(s.p)

	s.Assert().Equal(0, s.w.playerList.Count())
	s.Assert().Nil(s.w.playerRooms.playerToRoom[s.p])
	s.Assert().Equal(0, len(s.w.playerRooms.roomToPlayers.Get(s.w.StartRoom)))
	s.Assert().Equal(0, len(s.w.StartRoom.GetPlayers()))
}
