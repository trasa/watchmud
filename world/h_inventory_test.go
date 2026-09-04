package world

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"github.com/trasa/watchmud-message"
	"github.com/trasa/watchmud-message/slot"
	"github.com/trasa/watchmud/client"
	"github.com/trasa/watchmud/gameserver"
	"github.com/trasa/watchmud/object"
	"github.com/trasa/watchmud/player"
)

type handleInventorySuite struct {
	suite.Suite
	w *World
	r *player.Recorder
	p *player.Player
	c *client.TestClient
}

func TestHandleInventorySuite(t *testing.T) {
	suite.Run(t, new(handleInventorySuite))
}
func (s *handleInventorySuite) SetupTest() {
	s.w, _ = newTestWorld()
	s.r = &player.Recorder{}
	s.p = player.NewTestPlayer("foo", "foo", s.r)
	s.w.AddPlayer(s.p)
	s.c = client.NewTestClient(s.p)
}

func (s *handleInventorySuite) handleParameter() *gameserver.HandlerParameter {
	msg, err := message.NewGameMessage(message.InventoryRequest{})
	s.Assert().NoError(err)
	return gameserver.NewHandlerParameter(s.c, msg)
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
		InstanceId: uuid.New(),
		Definition: defnPtr,
	}
	s.Assert().NoError(s.p.Inventory().Add(instPtr))

	invHP := s.handleParameter()
	s.w.handleInventory(invHP)

	s.Assert().Equal(1, len(s.r.Sent))
	resp := s.r.Sent[0].(message.InventoryResponse)
	s.Assert().Equal(1, len(resp.InventoryItems))
	s.Assert().Equal(instPtr.Id(), resp.InventoryItems[0].Id)
	s.Assert().Equal(instPtr.Id(), resp.InventoryItems[0].Id)
}
