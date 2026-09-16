package player

import (
	"slices"
	"testing"
	"uuid"

	"github.com/stretchr/testify/suite"
)

type PlayersSuite struct {
	suite.Suite
	players *List
	p       *Player
	other   *Player
}

func TestPlayersSuite(t *testing.T) {
	suite.Run(t, new(PlayersSuite))
}

func (s *PlayersSuite) SetupTest() {
	s.players = NewList()
	s.p = NewTestPlayer(uuid.New(), "test", &Recorder{})
	s.other = NewTestPlayer(uuid.New(), "other", &Recorder{})
}

func (s *PlayersSuite) all() []*Player {
	return slices.Collect(s.players.All())
}

func (s *PlayersSuite) TestAdd() {
	s.players.Add(s.p)

	s.Assert().Same(s.p, s.players.FindByName("test"))
	s.Assert().Equal(1, s.players.Count())
}

func (s *PlayersSuite) TestRemove() {
	s.players.Add(s.p)
	s.players.Remove(s.p)

	s.Assert().Nil(s.players.FindByName("test"))
	s.Assert().Equal(0, s.players.Count())
}

// tolerated, and logged, rather than a panic in the middle of a move
func (s *PlayersSuite) TestRemove_DoesntExist() {
	s.players.Remove(s.p)

	s.Assert().Equal(0, s.players.Count())
}

func (s *PlayersSuite) TestAdd_Twice() {
	s.players.Add(s.p)
	s.players.Add(s.p)

	s.Assert().Equal([]*Player{s.p}, s.all())
}

func (s *PlayersSuite) TestFindByName_NotFound() {
	s.Assert().Nil(s.players.FindByName("nobody"))
}

// The order players are listed in is the order they joined, so a room
// description doesn't reshuffle between looks.
func (s *PlayersSuite) TestAllIsJoinOrder() {
	third := NewTestPlayer(uuid.New(), "third", &Recorder{})
	s.players.Add(s.p)
	s.players.Add(s.other)
	s.players.Add(third)

	s.Require().Equal([]*Player{s.p, s.other, third}, s.all())

	s.players.Remove(s.p)
	s.Assert().Equal([]*Player{s.other, third}, s.all())
}

func (s *PlayersSuite) TestAllExcept() {
	s.players.Add(s.p)
	s.players.Add(s.other)

	s.Assert().Equal([]*Player{s.other}, slices.Collect(s.players.AllExcept(s.p)))
}

// Nobody to exclude: the whole room hears it.
func (s *PlayersSuite) TestAllExcept_Nil() {
	s.players.Add(s.p)
	s.players.Add(s.other)

	s.Assert().Equal([]*Player{s.p, s.other}, slices.Collect(s.players.AllExcept(nil)))
}

// Slice is a copy: adding afterwards doesn't reach through it.
func (s *PlayersSuite) TestSliceIsACopy() {
	s.players.Add(s.p)

	all := s.players.Slice()
	s.players.Add(s.other)

	s.Assert().Equal([]*Player{s.p}, all)
}
