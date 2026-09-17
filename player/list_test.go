package player

import (
	"slices"
	"testing"
	"uuid"

	"github.com/stretchr/testify/suite"
)

type playerListSuite struct {
	suite.Suite
	list  *List
	p     *Player
	other *Player
}

func TestPlayerListSuite(t *testing.T) {
	suite.Run(t, new(playerListSuite))
}

func (s *playerListSuite) SetupTest() {
	s.list = NewList()
	s.p = NewTestPlayer(uuid.New(), "test", &Recorder{})
	s.other = NewTestPlayer(uuid.New(), "other", &Recorder{})
}

func (s *playerListSuite) all() []*Player {
	return slices.Collect(s.list.All())
}

func (s *playerListSuite) TestAdd() {
	s.list.Add(s.p)

	s.Assert().Same(s.p, s.list.FindByName("test"))
	s.Assert().Equal(1, s.list.Count())
}

func (s *playerListSuite) TestRemove() {
	s.list.Add(s.p)
	s.list.Remove(s.p)

	s.Assert().Nil(s.list.FindByName("test"))
	s.Assert().Equal(0, s.list.Count())
}

// tolerated, and logged, rather than a panic in the middle of a move
func (s *playerListSuite) TestRemove_DoesntExist() {
	s.list.Remove(s.p)

	s.Assert().Equal(0, s.list.Count())
}

func (s *playerListSuite) TestAdd_Twice() {
	s.list.Add(s.p)
	s.list.Add(s.p)

	s.Assert().Equal([]*Player{s.p}, s.all())
}

func (s *playerListSuite) TestFindByName_NotFound() {
	s.Assert().Nil(s.list.FindByName("nobody"))
}

// The order players are listed in is the order they joined, so a room
// description doesn't reshuffle between looks.
func (s *playerListSuite) TestAllIsJoinOrder() {
	third := NewTestPlayer(uuid.New(), "third", &Recorder{})
	s.list.Add(s.p)
	s.list.Add(s.other)
	s.list.Add(third)

	s.Require().Equal([]*Player{s.p, s.other, third}, s.all())

	s.list.Remove(s.p)
	s.Assert().Equal([]*Player{s.other, third}, s.all())
}

func (s *playerListSuite) TestAllExcept() {
	s.list.Add(s.p)
	s.list.Add(s.other)

	s.Assert().Equal([]*Player{s.other}, slices.Collect(s.list.AllExcept(s.p)))
}

// Nobody to exclude: the whole room hears it.
func (s *playerListSuite) TestAllExcept_Nil() {
	s.list.Add(s.p)
	s.list.Add(s.other)

	s.Assert().Equal([]*Player{s.p, s.other}, slices.Collect(s.list.AllExcept(nil)))
}

// Slice is a copy: adding afterwards doesn't reach through it.
func (s *playerListSuite) TestSliceIsACopy() {
	s.list.Add(s.p)

	all := s.list.Slice()
	s.list.Add(s.other)

	s.Assert().Equal([]*Player{s.p}, all)
}
