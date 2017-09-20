package world

import (
	"github.com/stretchr/testify/assert"
	"github.com/trasa/watchmud/message"
	"github.com/trasa/watchmud/object"
	"github.com/trasa/watchmud/player"
	"testing"
)

func TestInventory_success(t *testing.T) {
	w := newTestWorld()
	p := player.NewTestPlayer("guy")
	p.AddInventory(player.InventoryItem{
		Id:               "id",
		ShortDescription: "short desc",
		ObjectCategory:   object.TREASURE,
		Quantity:         3,
	})

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
	assert.Equal(t, "id", resp.InventoryItems[0].Id)
}
