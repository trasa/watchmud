package world

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/trasa/watchmud-message"
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
	s.other = player.NewTestPlayer("other", "other", nil)
	s.w.AddPlayer(s.other)
}

func (s *handleLookSuite) TestSuccess() {

	s.w.handleLook(s.handlerParameter(message.LookResponse{}))

	response := sent[message.LookResponse](s.T(), s.r, 0)
	s.Assert().True(response.Success)
	s.Assert().NotNil(response.RoomDescription.Name)
	s.Assert().NotNil(response.RoomDescription.Description)
	s.Assert().Equal(1, len(response.RoomDescription.Players))
	s.Assert().Equal("other", response.RoomDescription.Players[0])
}
