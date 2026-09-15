package world

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/trasa/watchmud/command"
	"github.com/trasa/watchmud/direction"
	"github.com/trasa/watchmud/event"
	"github.com/trasa/watchmud/gameserver"
)

type handleFleeSuite struct {
	worldTestSuite
	roller *loadedDice
}

func TestHandleFleeSuite(t *testing.T) {
	suite.Run(t, new(handleFleeSuite))
}

func (s *handleFleeSuite) SetupTest() {
	s.worldTestSuite.SetupTest()
	s.roller = newLoadedDice()
	s.w.roller = s.roller
}

func (s *handleFleeSuite) kill(target string) {
	s.T().Helper()
	cmd := command.Kill{Target: target}
	s.w.handleKill(gameserver.NewHandlerParameter(s.c, cmd), cmd)
}

func (s *handleFleeSuite) flee() {
	s.T().Helper()
	cmd := command.Flee{}
	s.w.handleFlee(gameserver.NewHandlerParameter(s.c, cmd), cmd)
}

func (s *handleFleeSuite) TestNotFighting() {
	s.flee()

	s.Assert().Equal(1, len(s.r.Sent))

	f := s.r.Sent[0].(event.Failed)
	s.Assert().Equal("flee", f.Verb)
	s.Assert().Equal(event.NoFight, f.Code)
}

func (s *handleFleeSuite) TestSuccess() {
	s.roller.Add(int(direction.North))
	s.kill("target")
	s.flee()
}
