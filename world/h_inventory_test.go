package world

import (
	"testing"
	"uuid"

	"github.com/stretchr/testify/suite"
	"github.com/trasa/watchmud/command"
	"github.com/trasa/watchmud/event"
	"github.com/trasa/watchmud/gameserver"
	"github.com/trasa/watchmud/object"
	"github.com/trasa/watchmud/player"
	"github.com/trasa/watchmud/slot"
)

type handleInventorySuite struct {
	suite.Suite
	w *World
	r *player.Recorder
	p *player.Player
	c *gameserver.TestConn
}

func TestHandleInventorySuite(t *testing.T) {
	suite.Run(t, new(handleInventorySuite))
}
func (s *handleInventorySuite) SetupTest() {
	s.w, _ = NewTestWorld()
	s.r = &player.Recorder{}
	s.p = player.NewTestPlayer(uuid.New(), "foo", s.r)
	s.w.AddPlayer(s.p)
	s.c = gameserver.NewTestConn(s.p)
}

func (s *handleInventorySuite) handleParameter() *gameserver.HandlerParameter {
	return gameserver.NewHandlerParameter(s.c, command.Inventory{})
}

func (s *handleInventorySuite) TestInventory_Success() {
	defnPtr := object.NewDefinition(
		"defnid",
		"name",
		"zone",
		object.Treasure,
		[]string{},
		"short desc",
		"in room",
		slot.None,
	)
	instPtr := &object.Instance{
		Id:         uuid.New(),
		Definition: defnPtr,
	}
	s.p.Inventory().Add(instPtr)

	s.w.handleInventory(s.handleParameter(), command.Inventory{})

	s.Assert().Equal(1, len(s.r.Sent))
	resp := s.r.Sent[0].(event.Inventory)
	s.Assert().Equal(1, len(resp.Items))
	s.Assert().Equal(instPtr.Id.String(), resp.Items[0].Id)
	s.Assert().Equal("short desc", resp.Items[0].ShortDescription)
}
