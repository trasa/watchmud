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
	s.Assert().Equal(s.w.VoidRoom, s.w.playerRoom(s.p))
	s.Assert().Equal(0, len(s.w.StartRoom.Players()))
}
