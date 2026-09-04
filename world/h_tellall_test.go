package world

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/trasa/watchmud-message"
	"github.com/trasa/watchmud/player"
)

type handleTellAllSuite struct {
	worldTestSuite
	receiver    *player.Player
	receiverRec *player.Recorder
	other       *player.Player
	otherRec    *player.Recorder
}

func TestHandleTellAllSuite(t *testing.T) {
	suite.Run(t, new(handleTellAllSuite))
}

func (s *handleTellAllSuite) SetupTest() {
	s.worldTestSuite.SetupTest()

	s.receiverRec = &player.Recorder{}
	s.receiver = player.NewTestPlayer("receiver", "receiver", s.receiverRec)
	s.w.AddPlayer(s.receiver)

	s.otherRec = &player.Recorder{}
	s.other = player.NewTestPlayer("other", "other", s.otherRec)
	s.w.AddPlayer(s.other)
}

func (s *handleTellAllSuite) TestSuccess() {

	s.w.handleTellAll(s.handlerParameter(message.TellAllRequest{Value: "hi"}))

	// did we tell otherPlayer?
	s.Assert().Equal(1, len(s.otherRec.Sent))
	s.Assert().Equal(1, len(s.receiverRec.Sent))

	// sender should have gotten response but NOT part of the send to all players
	s.Assert().Equal(1, len(s.r.Sent))
	senderResponse := sent[message.TellAllResponse](s.T(), s.r, 0)
	s.Assert().True(senderResponse.Success)
}

func (s *handleTellAllSuite) TestNoValue() {
	s.w.handleTellAll(s.handlerParameter(message.TellAllRequest{Value: ""}))

	// did we tell otherPlayer? (should be 0)
	s.Assert().Equal(0, len(s.otherRec.Sent))
	s.Assert().Equal(0, len(s.receiverRec.Sent))

	// sender should have gotten response but NOT part of the send to all players
	s.Assert().Equal(1, len(s.r.Sent))

	senderResponse := sent[message.TellAllResponse](s.T(), s.r, 0)
	s.Assert().False(senderResponse.Success)
}
