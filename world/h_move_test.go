package world

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/trasa/watchmud-message"
	"github.com/trasa/watchmud-message/direction"
	"github.com/trasa/watchmud/client"
	"github.com/trasa/watchmud/gameserver"
	"github.com/trasa/watchmud/player"
)

type HandleMoveSuite struct {
	suite.Suite
	w *World
	r *player.Recorder
	p *player.Player
	c *client.TestClient
}

func TestHandleMoveSuite(t *testing.T) {
	suite.Run(t, new(HandleMoveSuite))
}

func (s *HandleMoveSuite) SetupTest() {
	s.w, _ = newTestWorld()
	s.r = &player.Recorder{}
	s.p = player.NewTestPlayer("p", "p", s.r)
	s.w.AddPlayer(s.p)
	s.c = client.NewTestClient(s.p)
}

func (s *HandleMoveSuite) handlerParameter(dir direction.Direction) *gameserver.HandlerParameter {
	msg, err := message.NewGameMessage(message.MoveRequest{Direction: int32(dir)})
	s.Assert().NoError(err)
	return gameserver.NewHandlerParameter(s.c, msg)
}

func (s *HandleMoveSuite) TestMove_butYouCant() {
	s.w.handleMove(s.handlerParameter(direction.North))

	s.Assert().Equal(1, len(s.r.Sent))

	resp := s.r.Sent[0].(message.MoveResponse)
	s.Assert().False(resp.Success)
	s.Assert().Equal(resp.ResultCode, "CANT_GO_THAT_WAY")
}

func (s *HandleMoveSuite) TestMoveWhileFighting() {
	r := &player.Recorder{}
	other := player.NewTestPlayer("other", "other", r)
	s.w.AddPlayer(other)
	//s.w.fightLedger.Fight(s.p, other, s.w.StartRoom.Zone.Id, s.w.StartRoom.Id)

	//s.w.handleMove(s.handlerParameter(direction.North))

	//s.Assert().Equal(1, len(s.r.Sent))

	//resp := s.r.Sent[0].(message.MoveResponse)
	//s.Assert().False(resp.Success)
	//s.Assert().Equal(resp.ResultCode, "IN_A_FIGHT")
}
