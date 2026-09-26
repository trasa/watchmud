package world

import (
	"testing"
	"uuid"

	"github.com/stretchr/testify/suite"
	"github.com/watchmud/watchmud/command"
	"github.com/watchmud/watchmud/event"
	"github.com/watchmud/watchmud/gameserver"
	"github.com/watchmud/watchmud/player"
	"github.com/watchmud/watchmud/rules"
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

func (s *HandleMoveSuite) move(dir rules.Direction) {
	s.T().Helper()
	cmd := command.Move{Direction: dir}
	s.w.handleMove(gameserver.NewHandlerParameter(s.c, cmd), cmd)
}

func (s *HandleMoveSuite) TestMove_butYouCant() {
	s.move(rules.DirectionNorth)

	s.Assert().Equal(1, len(s.r.Sent))

	failed := s.r.Sent[0].(event.Failed)
	s.Assert().Equal("move", failed.Verb)
	s.Assert().Equal(event.CantGoThatWay, failed.Code)
}

func (s *HandleMoveSuite) TestMoveWhileFighting() {
	r := &player.Recorder{}
	other := player.NewTestPlayer(uuid.New(), "other", r)
	s.w.PlacePlayer(other, s.w.StartRoom)

	s.Assert().NoError(s.w.fightLedger.Fight(s.p, other, s.w.StartRoom.Zone.Id, s.w.StartRoom.Id))
	s.move(rules.DirectionNorth)

	failed := s.r.Sent[0].(event.Failed)
	s.Assert().Equal(event.InAFight, failed.Code)
}
