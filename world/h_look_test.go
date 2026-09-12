package world

import (
	"testing"
	"uuid"

	"github.com/stretchr/testify/suite"
	"github.com/trasa/watchmud/command"
	"github.com/trasa/watchmud/event"
	"github.com/trasa/watchmud/player"
)

type handleLookSuite struct {
	worldTestSuite
	other *player.Player
}

func TestHandleLookSuite(t *testing.T) {
	suite.Run(t, new(handleLookSuite))
}

func (s *handleLookSuite) SetupTest() {
	s.worldTestSuite.SetupTest()
	s.other = player.NewTestPlayer(uuid.New(), "other", nil)
	s.w.AddPlayer(s.other)
}

func (s *handleLookSuite) TestSuccess() {
	cmd := command.Look{}
	s.w.handleLook(s.handlerParameter(cmd), cmd)

	desc := sent[event.RoomDescription](s.T(), s.r, 0)
	s.Assert().NotEmpty(desc.Name)
	s.Assert().NotEmpty(desc.Description)
	s.Assert().Equal(1, len(desc.Players))
	s.Assert().Equal("other", desc.Players[0])
}
