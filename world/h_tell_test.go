package world

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/trasa/watchmud-message"
	"github.com/trasa/watchmud/client"
	"github.com/trasa/watchmud/gameserver"
	"github.com/trasa/watchmud/player"
)

type handleTellSuite struct {
	worldTestSuite
	sender       *player.Player
	senderRec    *player.Recorder
	senderClient *client.TestClient
	receiver     *player.Player
	receiverRec  *player.Recorder
}

func TestHandleTellSuite(t *testing.T) {
	suite.Run(t, new(handleTellSuite))
}

func (s *handleTellSuite) SetupTest() {
	s.worldTestSuite.SetupTest()
	s.senderRec = &player.Recorder{}
	s.sender = player.NewTestPlayer("sender", "sender", s.senderRec)
	s.w.AddPlayer(s.sender)
	s.senderClient = client.NewTestClient(s.sender)

	s.receiverRec = &player.Recorder{}
	s.receiver = player.NewTestPlayer("receiver", "receiver", s.receiverRec)
	s.w.AddPlayer(s.receiver)
}

func (s *handleTellSuite) handlerParameter(value string) *gameserver.HandlerParameter {
	msg, err := message.NewGameMessage(message.TellRequest{
		ReceiverPlayerName: s.receiver.Name,
		Value:              value,
	})
	s.Assert().NoError(err)
	return gameserver.NewHandlerParameter(s.senderClient, msg)
}

func (s *handleTellSuite) TestHandleTell() {
	s.w.handleTell(s.handlerParameter("hi"))

	// assert tell to receiver
	s.Assert().Equal(1, len(s.receiverRec.Sent))
	recdMessage := s.receiverRec.Sent[0].(message.TellNotification)
	s.Assert().Equal(s.sender.Name, recdMessage.Sender)
	s.Assert().Equal("hi", recdMessage.Value)

	// assert tell-response to sender
	s.Assert().Equal(1, len(s.senderRec.Sent))
	response := s.senderRec.Sent[0].(message.TellResponse)
	s.Assert().Equal("OK", response.ResultCode)
	s.Assert().True(response.Success)
}

func (s *handleTellSuite) ReceiverNotFound() {
	s.w.RemovePlayer(s.receiver)

	// act
	s.w.handleTell(s.handlerParameter("hi"))

	// assert tell-response to sender
	s.Assert().Equal(1, len(s.senderRec.Sent))
	response := s.senderRec.Sent[0].(message.TellResponse)
	s.Assert().Equal("TO_PLAYER_NOT_FOUND", response.ResultCode)
	s.Assert().False(response.Success)
}
