package world

import (
	"github.com/trasa/watchmud/message"
	"github.com/trasa/watchmud/object"
)

func (w *World) handleInventory(msg *message.IncomingMessage) {
	// TODO: something about getting the player's inventory here
	items := []message.InventoryItem{
		{
			Id:               "hi",
			ShortDescription: "hi there",
			ObjectCategory:   object.FOOD,
		},
	}
	resp := message.InventoryResponse{
		Response: message.NewSuccessfulResponse("inv"),
		InventoryItems: items,
	}
	msg.Player.Send(resp)
}