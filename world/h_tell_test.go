package world

import (
	"testing"
	"uuid"

	"github.com/stretchr/testify/suite"
	"github.com/trasa/watchmud/command"
	"github.com/trasa/watchmud/event"
	"github.com/trasa/watchmud/gameserver"
	"github.com/trasa/watchmud/player"
)

type handleTellSuite struct {
	worldTestSuite
	sender      *player.Player
	senderRec   *player.Recorder
	senderConn  *gameserver.TestConn
	receiver    *player.Player
	receiverRec *player.Recorder
}

func TestHandleTellSuite(t *testing.T) {
	suite.Run(t, new(handleTellSuite))
}

func (s *handleTellSuite) SetupTest() {
	s.worldTestSuite.SetupTest()
	s.senderRec = &player.Recorder{}
	s.sender = player.NewTestPlayer(uuid.New(), "sender", s.senderRec)
	s.w.AddPlayer(s.sender)
	s.senderConn = gameserver.NewTestConn(s.sender)

	s.receiverRec = &player.Recorder{}
	s.receiver = player.NewTestPlayer(uuid.New(), "receiver", s.receiverRec)
	s.w.AddPlayer(s.receiver)
}

func (s *handleTellSuite) tell(value string) {
	s.T().Helper()
	cmd := command.Tell{To: s.receiver.Name(), Value: value}
	s.w.handleTell(gameserver.NewHandlerParameter(s.senderConn, cmd), cmd)
}

func (s *handleTellSuite) TestHandleTell() {
	s.tell("hi")

	// one event, delivered to both ends of the conversation
	s.Assert().Equal(1, len(s.receiverRec.Sent))
	recd := s.receiverRec.Sent[0].(event.Told)
	s.Assert().Equal(s.sender.Name(), recd.From)
	s.Assert().Equal(s.receiver.Name(), recd.To)
	s.Assert().Equal("hi", recd.Value)

	s.Assert().Equal(1, len(s.senderRec.Sent))
	s.Assert().Equal(recd, s.senderRec.Sent[0].(event.Told))
}

func (s *handleTellSuite) ReceiverNotFound() {
	s.w.RemovePlayer(s.receiver)

	// act
	s.tell("hi")

	// assert failure to sender
	s.Assert().Equal(1, len(s.senderRec.Sent))
	failed := s.senderRec.Sent[0].(event.Failed)
	s.Assert().Equal("tell", failed.Verb)
	s.Assert().Equal(event.ToPlayerNotFound, failed.Code)
}
