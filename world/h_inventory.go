package world

import (
	"github.com/trasa/watchmud-message"
	"github.com/trasa/watchmud/gameserver"
)

func (w *World) handleInventory(msg *gameserver.HandlerParameter) {
	var items []*message.InventoryResponse_InventoryItem
	for _, instPtr := range msg.Player.Inventory().GetAll() {
		if !msg.Player.Slots().IsItemInUse(instPtr) {
			items = append(items, &message.InventoryResponse_InventoryItem{
				Id:               instPtr.Id(),
				ShortDescription: instPtr.Definition.ShortDescription,
				ObjectCategories: instPtr.Definition.Categories.ToInt32List(),
			})
		}
	}
	msg.Player.Send(message.InventoryResponse{
		Success:        true,
		ResultCode:     "OK",
		InventoryItems: items,
	})
}
