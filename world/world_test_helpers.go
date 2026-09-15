package world

import (
	"os"
	"testing"
	"uuid"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/trasa/watchmud/command"
	"github.com/trasa/watchmud/gameserver"
	"github.com/trasa/watchmud/loader"
	"github.com/trasa/watchmud/memstore"
	"github.com/trasa/watchmud/player"
)

type worldTestSuite struct {
	suite.Suite
	w *World
	r *player.Recorder
	p *player.Player
	c *gameserver.TestConn
}

// sent returns the ith message in the recorder
func sent[T any](t *testing.T, r *player.Recorder, i int) T {
	t.Helper()
	require.Greater(t, len(r.Sent), i, "wanted Sent[%d], only %d sent", i, len(r.Sent))
	v, ok := r.Sent[i].(T)
	require.Truef(t, ok, "Sent[%d] is %T, want %T", i, r.Sent[i], *new(T))
	return v
}

// handlerParameter builds the handler input for a command.
func (s *worldTestSuite) handlerParameter(cmd command.Command) *gameserver.HandlerParameter {
	s.T().Helper()
	return gameserver.NewHandlerParameter(s.c, cmd)
}

func (s *worldTestSuite) SetupTest() {
	w, err := NewTestWorld()
	require.NoError(s.T(), err)
	s.w = w
	s.r = &player.Recorder{}
	s.p = player.NewTestPlayer(uuid.New(), "testdood", s.r)
	s.w.AddPlayer(s.p)
	s.c = gameserver.NewTestConn(s.p)
}

func NewTestWorld() (*World, error) {
	fs := os.DirFS("../testcontent")
	content, err := loader.LoadContent(fs)
	if err != nil {
		return nil, err
	}

	store := memstore.New()
	roller := newLoadedDice()

	return New(content, store, roller)
}
