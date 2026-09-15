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

type handleFleeSuite struct {
	worldTestSuite
	roller   *loadedDice
	other    *player.Player
	otherRec *player.Recorder
}

func TestHandleFleeSuite(t *testing.T) {
	suite.Run(t, new(handleFleeSuite))
}

func (s *handleFleeSuite) SetupTest() {
	s.worldTestSuite.SetupTest()
	s.roller = newLoadedDice()
	s.w.roller = s.roller
	s.otherRec = &player.Recorder{}
	s.other = player.NewTestPlayer(uuid.New(), "other", s.otherRec)
	s.w.AddPlayer(s.other)
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
	s.roller.Add(int(direction.South))
	s.kill("target")
	s.flee()

	// attack
	attacking := s.r.Sent[0].(event.Attacking)
	s.Assert().Equal("Target Drone", attacking.Target)

	// flee
	fleeing := s.r.Sent[1].(event.Fleeing)
	s.Assert().Equal(s.p.Name(), fleeing.Who)

	// success
	fled := s.r.Sent[2].(event.Fled)
	s.Assert().Equal(s.p.Name(), fled.Who)

	// what did the other see?
	s.Assert().Equal(s.p.Name(), s.otherRec.Sent[0].(event.Fleeing).Who)
	s.Assert().Equal(s.p.Name(), s.otherRec.Sent[1].(event.Fled).Who)
	s.Assert().Equal(direction.South, s.otherRec.Sent[2].(event.Left).Direction)
}

func (s *handleFleeSuite) TestNoEscape() {
	rolls := []int{int(direction.North),
		int(direction.North),
		int(direction.North),
		int(direction.North),
		int(direction.North),
		int(direction.North),
	}
	s.roller.Load(rolls)

	s.kill("target")
	s.flee()

	// attack
	attacking := s.r.Sent[0].(event.Attacking)
	s.Assert().Equal("Target Drone", attacking.Target)

	// flee
	fleeing := s.r.Sent[1].(event.Fleeing)
	s.Assert().Equal(s.p.Name(), fleeing.Who)

	// success
	fled := s.r.Sent[2].(event.FleeAttemptFailed)
	s.Assert().Equal(s.p.Name(), fled.Who)
}
