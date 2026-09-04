package world

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/trasa/watchmud-message"
	"github.com/trasa/watchmud/client"
	"github.com/trasa/watchmud/gameserver"
	"github.com/trasa/watchmud/player"
)

type HandleKillSuite struct {
	suite.Suite
	w *World
	r *player.Recorder
	p *player.Player
	c *client.TestClient
}

func TestHandleKillSuite(t *testing.T) {
	suite.Run(t, new(HandleKillSuite))
}

func (s *HandleKillSuite) SetupTest() {
	s.w, _ = newTestWorld()
	s.r = &player.Recorder{}
	s.p = player.NewTestPlayer("testdood", "testdood", s.r)
	s.w.AddPlayer(s.p)
	s.c = client.NewTestClient(s.p)
}

func (s *HandleKillSuite) handleParameter(target string) *gameserver.HandlerParameter {
	msg, err := message.NewGameMessage(message.KillRequest{Target: target})
	s.Assert().NoError(err)
	return gameserver.NewHandlerParameter(s.c, msg)
}

func (s *HandleKillSuite) TestSuccess() {
	_ /*mob*/, _ = s.w.StartRoom.FindMobile("target")
	killHP := s.handleParameter("target")

	s.w.handleKill(killHP)

	s.Assert().Equal(1, len(s.r.Sent))
	resp := s.r.Sent[0].(message.KillResponse)
	s.Assert().True(resp.Success)
	s.Assert().Equal("OK", resp.ResultCode)

	//s.Assert().True(s.w.fightLedger.IsFighting(s.p))
	//s.Assert().Equal(mob, s.w.fightLedger.GetFight(s.p).Fightee)
	//s.Assert().Equal(s.p, s.w.fightLedger.GetFight(s.p).Fighter)
}

/*
func (s *HandleKillSuite) TestAlreadyFighting() {
	mob, _ := s.w.StartRoom.FindMobile("target")
	s.w.fightLedger.Fight(s.p, mob, s.w.StartRoom.Zone.Id, s.w.StartRoom.Id)

	killHP := newKillRequestHandleParameter(s.T(), s.c, "targetMob")

	s.w.handleKill(killHP)

	s.Assert().Equal(1, s.p.SentMessageCount())
	resp := s.p.GetSentResponse(0).(message.KillResponse)
	s.Assert().False(resp.Success)
	s.Assert().Equal("ALREADY_FIGHTING", resp.ResultCode)
}
*/
/*
func (s *HandleKillSuite) TestNoTarget() {
	killHP := newKillRequestHandleParameter(s.T(), s.c, "targetMob")

	s.w.handleKill(killHP)

	s.Assert().Equal(1, s.p.SentMessageCount())
	resp := s.p.GetSentResponse(0).(message.KillResponse)

	s.Assert().False(resp.Success)
	s.Assert().Equal("TARGET_NOT_FOUND", resp.ResultCode)
}

func (s *HandleKillSuite) TestNoFight() {
	mob, _ := s.w.StartRoom.FindMobile("target")
	mob.Definition.SetFlag(mobile.FlagNoFight)
	killHP := newKillRequestHandleParameter(s.T(), s.c, "target")

	s.w.handleKill(killHP)

	s.Assert().Equal(1, s.p.SentMessageCount())
	resp := s.p.GetSentResponse(0).(message.KillResponse)

	s.Assert().False(resp.Success)
	s.Assert().Equal("NO_FIGHT", resp.ResultCode)
}

func (s *HandleKillSuite) TestNoFightInRoom() {

	s.w.StartRoom.SetFlag(spaces.RoomFlagNoFight)

	killHP := newKillRequestHandleParameter(s.T(), s.c, "target")

	s.w.handleKill(killHP)

	s.Assert().Equal(1, s.p.SentMessageCount())
	resp := s.p.GetSentResponse(0).(message.KillResponse)

	s.Assert().False(resp.Success)
	s.Assert().Equal("NO_FIGHT_ROOM", resp.ResultCode)
}
*/
