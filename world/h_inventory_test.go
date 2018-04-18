package world

import (
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/trasa/watchmud/client"
	"github.com/trasa/watchmud/gameserver"
	"github.com/trasa/watchmud/message"
	"github.com/trasa/watchmud/object"
	"github.com/trasa/watchmud/player"
	"github.com/trasa/watchmud/slot"
	"testing"
)

func newInventoryRequestHandlerParameter(t *testing.T, c *client.TestClient) *gameserver.HandlerParameter {
	msg, err := message.NewGameMessage(message.InventoryRequest{})
	assert.NoError(t, err)
	return gameserver.NewHandlerParameter(c, msg)
}

func TestInventory_success(t *testing.T) {
	w := newTestWorld()
	p := player.NewTestPlayer("guy")
	c := client.NewTestClient(p)

	defnPtr := object.NewDefinition("defnid", "name", "zone",
		object.Treasure, []string{}, "short desc", "in room", slot.None)
	instPtr := &object.Instance{
		InstanceId: uuid.Must(uuid.NewV4()),
		Definition: defnPtr,
	}
	p.AddInventory(instPtr)

	invHP := newInventoryRequestHandlerParameter(t, c)
	w.handleInventory(invHP)

	assert.Equal(t, 1, p.SentMessageCount())
	resp := p.GetSentResponse(0).(message.InventoryResponse)
	assert.Equal(t, 1, len(resp.InventoryItems))
	assert.Equal(t, instPtr.Id(), resp.InventoryItems[0].Id)
}
