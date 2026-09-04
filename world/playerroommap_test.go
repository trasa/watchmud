package world

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/trasa/watchmud/player"
	"github.com/trasa/watchmud/spaces"
)

type PlayerRoomMapSuite struct {
	suite.Suite
	m         *PlayerRoomMap
	bob       *player.Player
	alice     *player.Player
	northRoom *spaces.Room
}

func TestPlayerRoomMapSuite(t *testing.T) {
	suite.Run(t, new(PlayerRoomMapSuite))
}

func (s *PlayerRoomMapSuite) SetupTest() {
	s.m = NewPlayerRoomMap()
	s.bob = player.NewTestPlayer("bob", "bob", &player.Recorder{})
	s.alice = player.NewTestPlayer("alice", "alice", &player.Recorder{})
	s.northRoom = spaces.NewTestRoom("north")
	s.m.Add(s.bob, s.northRoom)
	s.m.Add(s.alice, s.northRoom)
}

func (s *PlayerRoomMapSuite) TestGetPlayers() {
	s.Assert().Equal(2, len(s.m.GetPlayers(s.northRoom)))
	s.m.Remove(s.bob)
	s.Assert().Equal(1, len(s.m.GetPlayers(s.northRoom)))
}
