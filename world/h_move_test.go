package world

import (
	"testing"
	"uuid"

	"github.com/stretchr/testify/suite"
	"github.com/trasa/watchmud/command"
	"github.com/trasa/watchmud/direction"
	"github.com/trasa/watchmud/event"
	"github.com/trasa/watchmud/gameserver"
	"github.com/trasa/watchmud/player"
)

type HandleMoveSuite struct {
	worldTestSuite
}

func TestHandleMoveSuite(t *testing.T) {
	suite.Run(t, new(HandleMoveSuite))
}

func (s *HandleMoveSuite) SetupTest() {
	s.worldTestSuite.SetupTest()
}

func (s *HandleMoveSuite) move(dir direction.Direction) {
	s.T().Helper()
	cmd := command.Move{Direction: dir}
	s.w.handleMove(gameserver.NewHandlerParameter(s.c, cmd), cmd)
}

func (s *HandleMoveSuite) TestMove_butYouCant() {
	s.move(direction.North)

	s.Assert().Equal(1, len(s.r.Sent))

	failed := s.r.Sent[0].(event.Failed)
	s.Assert().Equal("move", failed.Verb)
	s.Assert().Equal(event.CantGoThatWay, failed.Code)
}

func (s *HandleMoveSuite) TestMoveWhileFighting() {
	r := &player.Recorder{}
	other := player.NewTestPlayer(uuid.New(), "other", r)
	s.w.AddPlayer(other)

	s.Assert().NoError(s.w.fightLedger.Fight(s.p, other, s.w.StartRoom.Zone.Id, s.w.StartRoom.Id))
	s.move(direction.North)

	failed := s.r.Sent[0].(event.Failed)
	s.Assert().Equal(event.InAFight, failed.Code)
}
