package world

import (
	"testing"
	"uuid"

	"github.com/stretchr/testify/suite"
	"github.com/watchmud/watchmud/command"
	"github.com/watchmud/watchmud/event"
	"github.com/watchmud/watchmud/gameserver"
	"github.com/watchmud/watchmud/player"
	"github.com/watchmud/watchmud/rules"
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
	s.w.PlacePlayer(s.p, s.w.StartRoom)
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
	s.equip(command.Equip{Slot: rules.SlotWield})

	s.Assert().Equal(1, len(s.r.Sent))
	failed := s.r.Sent[0].(event.Failed)
	s.Assert().Equal(event.NoTarget, failed.Code)
}
