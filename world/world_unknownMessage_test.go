package world

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/trasa/watchmud/event"
	"github.com/trasa/watchmud/gameserver"
)

type worldUnknownMessageSuite struct {
	worldTestSuite
}

func TestUnknownMessageSuite(t *testing.T) {
	suite.Run(t, new(worldUnknownMessageSuite))
}
func (s *worldUnknownMessageSuite) SetupTest() {
	s.worldTestSuite.SetupTest()
}

// unknownCommand is a command with no case in the dispatch switch. The
// compiler can't catch that -- a type switch has no exhaustiveness check --
// so the default arm has to behave.
type unknownCommand struct{}

func (unknownCommand) Verb() string { return "florb" }

func (s *worldUnknownMessageSuite) TestUnknownCommand() {
	h := gameserver.NewHandlerParameter(s.c, unknownCommand{})

	err := s.w.HandleIncomingMessage(h)

	s.Assert().Error(err)
	failed := sent[event.Failed](s.T(), s.r, 0)
	s.Assert().Equal("florb", failed.Verb)
	s.Assert().Equal(event.UnknownCommand, failed.Code)
}
