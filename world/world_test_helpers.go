package world

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	message "github.com/trasa/watchmud-message"
	"github.com/trasa/watchmud/client"
	"github.com/trasa/watchmud/gameserver"
	"github.com/trasa/watchmud/player"
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
