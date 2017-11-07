package world

import (
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/trasa/watchmud/message"
	"github.com/trasa/watchmud/object"
	"github.com/trasa/watchmud/player"
	"testing"
)

func TestInventory_success(t *testing.T) {
	w := newTestWorld()
	p := player.NewTestPlayer("guy")

	defnPtr := object.NewDefinition("defnid", "name", "zone",
		object.TREASURE, []string{}, "short desc", "in room")
	instPtr := &object.Instance{
		InstanceId: uuid.NewV4(),
		Definition: defnPtr,
	}

	p.AddInventory(instPtr)

	msg := message.IncomingMessage{
		Player: p,
		Request: message.InventoryRequest{
			Request: message.RequestBase{MessageType: "inv"},
		},
	}
	w.handleInventory(&msg)

	assert.Equal(t, 1, p.SentMessageCount())
	resp := p.GetSentResponse(0).(message.InventoryResponse)
	assert.Equal(t, 1, len(resp.InventoryItems))
	assert.Equal(t, instPtr.Id(), resp.InventoryItems[0].Id)
}
