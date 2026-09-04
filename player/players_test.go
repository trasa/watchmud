package player

import (
	"log"
	"testing"

	"github.com/stretchr/testify/suite"
)

type PlayersSuite struct {
	suite.Suite
	players *List
	p       *Player
	r       *Recorder
}

func TestPlayersSuite(t *testing.T) {
	suite.Run(t, new(PlayersSuite))
}

func (s *PlayersSuite) SetupTest() {
	s.players = NewList()
	s.r = &Recorder{}
	s.p = NewTestPlayer("test", "test", s.r)
}

func (s *PlayersSuite) TestAdd() {
	s.players.Add(s.p)

	_, ok := s.players.players[s.p]
	s.Assert().True(ok)
}

func (s *PlayersSuite) TestRemove() {
	s.players.Add(s.p)
	s.players.Remove(s.p)

	_, ok := s.players.players[s.p]
	s.Assert().False(ok) // not found
}

func (s *PlayersSuite) TestRemove_DoesntExist() {
	s.players.Remove(s.p)
	_, ok := s.players.players[s.p]
	s.Assert().False(ok)
}

func (s *PlayersSuite) TestIter() {
	s.players.Add(s.p)
	count := 0

	s.players.Iter(func(p *Player) {
		count++
	})

	s.Assert().Equal(1, count)
}

func (s *PlayersSuite) TestGetAll() {
	s.players.Add(s.p)
	other := NewTestPlayer("other", "other", &Recorder{})

	all := s.players.GetAll()
	s.players.Add(other)
	log.Printf("addr of all: %p", &all)

	s.Assert().Equal(1, len(all))
}
