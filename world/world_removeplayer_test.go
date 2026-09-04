package world

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type worldRemovePlayerSuite struct {
	worldTestSuite
}

func TestRemovePlayerSuite(t *testing.T) {
	suite.Run(t, new(worldRemovePlayerSuite))
}

func (s *worldRemovePlayerSuite) SetupTest() {
	s.worldTestSuite.SetupTest()
}

func (s *worldRemovePlayerSuite) TestRemovePlayer() {
	s.w.RemovePlayer(s.p)

	s.Assert().Equal(0, s.w.playerList.Count())
	s.Assert().Nil(s.w.playerRooms.playerToRoom[s.p])
	s.Assert().Equal(0, len(s.w.playerRooms.roomToPlayers.Get(s.w.StartRoom)))
	s.Assert().Equal(0, len(s.w.StartRoom.GetPlayers()))
}
