package world

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/trasa/watchmud-message"
	"github.com/trasa/watchmud/gameserver"
	"github.com/trasa/watchmud/player"
)

type handleLookSuite struct {
	worldTestSuite
	other *player.Player
}

func TestHandleLookSuite(t *testing.T) {
	suite.Run(t, new(handleLookSuite))
}

func (s *handleLookSuite) SetupTest() {
	s.worldTestSuite.SetupTest()
	s.other = player.NewTestPlayer("other", "other", &player.Recorder{})
	s.w.AddPlayer(s.other)
}

func (s *handleLookSuite) handlerParameter() *gameserver.HandlerParameter {
	msg, err := message.NewGameMessage(message.LookRequest{})
	s.Assert().NoError(err)
	return gameserver.NewHandlerParameter(s.c, msg)
}

func (s *handleLookSuite) TestLook_Successful() {

	s.w.HandleIncomingMessage(s.handlerParameter())

	resp := s.r.Sent[0].(message.LookResponse)
	s.Assert().True(resp.Success)
	s.Assert().NotNil(resp.RoomDescription.Name)
	s.Assert().NotNil(resp.RoomDescription.Description)
	s.Assert().Equal(1, len(resp.RoomDescription.Players))
	s.Assert().Equal("other", resp.RoomDescription.Players[0])
}
