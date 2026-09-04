package world

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/trasa/watchmud-message"
	"github.com/trasa/watchmud/player"
)

type WhoSuite struct {
	worldTestSuite
}

func TestWhoSuite(t *testing.T) {
	suite.Run(t, new(WhoSuite))
}

func (s *WhoSuite) SetupTest() {
	s.worldTestSuite.SetupTest()
}

func (s *WhoSuite) TestSuccess() {
	s.w.handleWho(s.handlerParameter(message.WhoRequest{}))

	resp := sent[message.WhoResponse](s.T(), s.r, 0)
	s.Assert().True(resp.Success)
	s.Assert().Equal(1, len(resp.PlayerInfo))
	s.Assert().Equal("testdood", resp.PlayerInfo[0].PlayerName)
	s.Assert().NotEqual("", resp.PlayerInfo[0].ZoneName)
	s.Assert().NotEqual("", resp.PlayerInfo[0].RoomName)
}

func (s *WhoSuite) TestNotInRoom() {
	s.w.playerRooms.Remove(s.p)

	s.w.handleWho(s.handlerParameter(message.WhoRequest{}))

	resp := sent[message.WhoResponse](s.T(), s.r, 0)
	s.Assert().True(resp.Success)
	s.Assert().Equal("", resp.PlayerInfo[0].ZoneName)
	s.Assert().Equal("", resp.PlayerInfo[0].RoomName)
}

func (s *WhoSuite) TestSort() {
	rec := &player.Recorder{}
	otherPlayer := player.NewTestPlayer("other", "other", rec)
	s.w.AddPlayer(otherPlayer)

	s.w.handleWho(s.handlerParameter(message.WhoRequest{}))
	response := sent[message.WhoResponse](s.T(), s.r, 0)

	s.Assert().Equal("testdood", response.PlayerInfo[0].PlayerName)
	s.Assert().Equal("other", response.PlayerInfo[1].PlayerName)
}

func (s *WhoSuite) TestLogoutRemovesPlayer() {

	rec := &player.Recorder{}
	otherPlayer := player.NewTestPlayer("other", "other", rec)
	s.w.AddPlayer(otherPlayer)
	s.w.RemovePlayer(otherPlayer)

	s.w.handleWho(s.handlerParameter(message.WhoRequest{}))

	response := sent[message.WhoResponse](s.T(), s.r, 0)
	s.Assert().Equal("testdood", response.PlayerInfo[0].PlayerName)
}
