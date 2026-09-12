package world

import (
	"testing"
	"uuid"

	"github.com/stretchr/testify/suite"
	"github.com/trasa/watchmud/command"
	"github.com/trasa/watchmud/event"
	"github.com/trasa/watchmud/gameserver"
	"github.com/trasa/watchmud/player"
	"github.com/trasa/watchmud/slot"
)

type HandleEquipSuite struct {
	suite.Suite
	w *World
	r *player.Recorder
	p *player.Player
	c *gameserver.TestConn
}

func TestHandleEquipSuite(t *testing.T) {
	suite.Run(t, new(HandleEquipSuite))
}

func (s *HandleEquipSuite) SetupTest() {
	s.w, _ = NewTestWorld()
	s.r = &player.Recorder{}
	s.p = player.NewTestPlayer(uuid.New(), "foo", s.r)
	s.w.AddPlayer(s.p)
	s.c = gameserver.NewTestConn(s.p)
}

func (s *HandleEquipSuite) equip(cmd command.Equip) {
	s.T().Helper()
	s.w.handleEquip(gameserver.NewHandlerParameter(s.c, cmd), cmd)
}

func (s *HandleEquipSuite) TestNoSlot() {
	s.equip(command.Equip{})

	s.Assert().Equal(1, len(s.r.Sent))
	failed := s.r.Sent[0].(event.Failed)
	s.Assert().Equal("equip", failed.Verb)
	s.Assert().Equal(event.NoSlotGiven, failed.Code)
}

func (s *HandleEquipSuite) TestNoTarget() {
	s.equip(command.Equip{Slot: slot.Wield})

	s.Assert().Equal(1, len(s.r.Sent))
	failed := s.r.Sent[0].(event.Failed)
	s.Assert().Equal(event.NoTarget, failed.Code)
}
