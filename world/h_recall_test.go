package world

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/watchmud/watchmud/command"
	"github.com/watchmud/watchmud/event"
)

type handleRecallSuite struct {
	worldTestSuite
}

func TestHandleRecallSuite(t *testing.T) {
	suite.Run(t, new(handleRecallSuite))
}

func (s *handleRecallSuite) SetupTest() {
	s.worldTestSuite.SetupTest()
	market, found := s.w.findRoomById("wrathrock", "market_square")
	s.Require().True(found)
	s.w.movePlayerMagically(s.p, market)
	s.r.Clear()
}

func (s *handleRecallSuite) recall() {
	s.T().Helper()
	s.Require().NoError(s.w.HandleIncomingMessage(s.handlerParameter(command.Recall{})))
}

func (s *handleRecallSuite) TestTakesYouToTheStart() {
	s.recall()

	s.Assert().Same(s.w.StartRoom, s.w.playerRoom(s.p))
	s.Assert().Equal(s.w.StartRoom.Name, sent[event.RoomDescription](s.T(), s.r, 0).Name)
}

// Recall used to be a free, certain way out of any fight -- which made flee,
// with its chance of failing, pointless. It is refused, the way walking off is.
func (s *handleRecallSuite) TestNotInAFight() {
	drone, found := s.w.StartRoom.FindMobile("target")
	s.Require().True(found)
	s.Require().NoError(s.w.fightLedger.Fight(s.p, drone, s.w.StartRoom.Zone.Id, s.w.StartRoom.Id))
	market := s.w.playerRoom(s.p)

	s.recall()

	s.Assert().Same(market, s.w.playerRoom(s.p), "still where they were")
	s.Assert().Equal([]any{event.Failed{Verb: "recall", Code: event.InAFight}}, s.r.Sent)
}
