package world

import (
	"testing"
	"uuid"

	"github.com/stretchr/testify/suite"
	"github.com/trasa/watchmud/command"
	"github.com/trasa/watchmud/event"
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
	s.receiver = player.NewTestPlayer(uuid.New(), "receiver", s.receiverRec)
	s.w.AddPlayer(s.receiver)

	s.otherRec = &player.Recorder{}
	s.other = player.NewTestPlayer(uuid.New(), "other", s.otherRec)
	s.w.AddPlayer(s.other)
}

func (s *handleTellAllSuite) TestSuccess() {

	cmd := command.TellAll{Value: "hi"}
	s.w.handleTellAll(s.handlerParameter(cmd), cmd)

	// did we tell otherPlayer?
	s.Assert().Equal(1, len(s.otherRec.Sent))
	s.Assert().Equal(1, len(s.receiverRec.Sent))

	// sender should have gotten response but NOT part of the send to all players
	s.Assert().Equal(1, len(s.r.Sent))
	shouted := sent[event.Shouted](s.T(), s.r, 0)
	s.Assert().Equal("testdood", shouted.Speaker)
	s.Assert().Equal("hi", shouted.Value)
}

func (s *handleTellAllSuite) TestNoValue() {
	cmd := command.TellAll{Value: ""}
	s.w.handleTellAll(s.handlerParameter(cmd), cmd)

	// did we tell otherPlayer? (should be 0)
	s.Assert().Equal(0, len(s.otherRec.Sent))
	s.Assert().Equal(0, len(s.receiverRec.Sent))

	// sender should have gotten response but NOT part of the send to all players
	s.Assert().Equal(1, len(s.r.Sent))

	failed := sent[event.Failed](s.T(), s.r, 0)
	s.Assert().Equal(event.NoValue, failed.Code)
}
