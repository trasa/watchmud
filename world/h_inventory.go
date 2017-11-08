package world

import (
	"github.com/trasa/watchmud/gameserver"
	"github.com/trasa/watchmud/message"
)

func (w *World) handleInventory(msg *gameserver.HandlerParameter) {
	items := []*message.InventoryResponse_InventoryItem{}
	for _, instPtr := range msg.Player.GetAllInventory() {
		items = append(items, &message.InventoryResponse_InventoryItem{
			Id:               instPtr.Id(),
			ShortDescription: instPtr.Definition.ShortDescription,
			ObjectCategories: instPtr.Definition.Categories.ToInt32List(),
		})
	}
	msg.Player.Send(message.InventoryResponse{
		Success:        true,
		ResultCode:     "OK",
		InventoryItems: items,
	})
}
