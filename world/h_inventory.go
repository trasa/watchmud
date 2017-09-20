package world

import (
	"github.com/trasa/watchmud/message"
	"github.com/trasa/watchmud/player"
)

func (w *World) handleInventory(msg *message.IncomingMessage) {
	items := []message.InventoryDescription{}
	for _, vals := range msg.Player.GetInventory() {
		for _, val := range vals {
			items = append(items, InventoryItemToDescription(val))
		}
	}
	resp := message.InventoryResponse{
		Response:       message.NewSuccessfulResponse("inv"),
		InventoryItems: items,
	}
	msg.Player.Send(resp)
}

func InventoryItemToDescription(item player.InventoryItem) message.InventoryDescription {
	return message.InventoryDescription{
		Id:               item.Id,
		ShortDescription: item.ShortDescription,
		ObjectCategory:   item.ObjectCategory,
		Quantity:         item.Quantity,
	}
}
