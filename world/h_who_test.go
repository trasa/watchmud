package world

import (
	"testing"
	"uuid"

	"github.com/stretchr/testify/suite"
	"github.com/watchmud/watchmud/command"
	"github.com/watchmud/watchmud/event"
	"github.com/watchmud/watchmud/player"
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
	s.w.handleWho(s.handlerParameter(command.Who{}), command.Who{})

	resp := sent[event.Who](s.T(), s.r, 0)
	s.Assert().Equal(1, len(resp.Players))
	s.Assert().Equal("testdood", resp.Players[0].PlayerName)
	s.Assert().NotEqual("", resp.Players[0].ZoneName)
	s.Assert().NotEqual("", resp.Players[0].RoomName)
}

func (s *WhoSuite) TestSort() {
	rec := &player.Recorder{}
	otherPlayer := player.NewTestPlayer(uuid.New(), "other", rec)
	s.w.PlacePlayer(otherPlayer, s.w.StartRoom)

	s.w.handleWho(s.handlerParameter(command.Who{}), command.Who{})
	response := sent[event.Who](s.T(), s.r, 0)

	// TODO what sort order is the command working with?
	s.Assert().Equal("other", response.Players[0].PlayerName)
	s.Assert().Equal("testdood", response.Players[1].PlayerName)
}

func (s *WhoSuite) TestLogoutRemovesPlayer() {

	rec := &player.Recorder{}
	otherPlayer := player.NewTestPlayer(uuid.New(), "other", rec)
	s.w.PlacePlayer(otherPlayer, s.w.StartRoom)
	s.w.RemovePlayer(otherPlayer)

	s.w.handleWho(s.handlerParameter(command.Who{}), command.Who{})

	response := sent[event.Who](s.T(), s.r, 0)
	s.Assert().Equal("testdood", response.Players[0].PlayerName)
}
