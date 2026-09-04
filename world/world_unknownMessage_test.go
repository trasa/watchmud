package world

import (
	"testing"

	"github.com/stretchr/testify/suite"
	message "github.com/trasa/watchmud-message"
	"github.com/trasa/watchmud/gameserver"
)

type worldUnknownMessageSuite struct {
	worldTestSuite
}

func TestUnknownMessageSuite(t *testing.T) {
	suite.Run(t, new(worldUnknownMessageSuite))
}
func (s *worldUnknownMessageSuite) SetupTest() {
	s.worldTestSuite.SetupTest()
}

func (s *worldUnknownMessageSuite) TestUnknownMessageType() {
	m := &message.GameMessage{} // not a valid message, has no inner type
	h := gameserver.NewHandlerParameter(s.c, m)

	s.w.HandleIncomingMessage(h)

	resp := sent[message.ErrorResponse](s.T(), s.r, 0)
	s.Assert().False(resp.Success)
	s.Assert().Equal("UNKNOWN_MESSAGE_TYPE", resp.ResultCode)
}
