package world

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/trasa/watchmud-message"
	"github.com/trasa/watchmud/client"
	"github.com/trasa/watchmud/gameserver"
	"github.com/trasa/watchmud/player"
)

type HandleEquipSuite struct {
	suite.Suite
	w      *World
	r      *player.Recorder
	p      *player.Player
	c      *client.TestClient
	msg    *message.GameMessage
	handle *gameserver.HandlerParameter
}

func TestHandleEquipSuite(t *testing.T) {
	suite.Run(t, new(HandleEquipSuite))
}

func (s *HandleEquipSuite) SetupTest() {
	s.w, _ = newTestWorld()
	s.r = &player.Recorder{}
	s.p = player.NewTestPlayer("foo", "foo", s.r)
	s.w.AddPlayer(s.p)
	s.c = client.NewTestClient(s.p)

	msg, err := message.NewGameMessage(message.EquipRequest{})
	s.Assert().NoError(err)
	s.msg = msg
	s.handle = gameserver.NewHandlerParameter(s.c, s.msg)
}

func (s *HandleEquipSuite) TestNoSlot() {
	s.w.handleEquip(s.handle)

	s.Assert().Equal(1, len(s.r.Sent))
	resp := s.r.Sent[0].(message.EquipResponse)
	s.Assert().False(resp.Success)
	s.Assert().Equal("NO_SLOT_GIVEN", resp.ResultCode)
}

func (s *HandleEquipSuite) TestNoTarget() {
	s.msg.GetEquipRequest().SlotLocation = 1
	s.w.handleEquip(s.handle)

	s.Assert().Equal(1, len(s.r.Sent))
	resp := s.r.Sent[0].(message.EquipResponse)
	s.Assert().False(resp.Success)
	s.Assert().Equal("NO_TARGET", resp.ResultCode)
}
