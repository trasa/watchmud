package world

import (
	"github.com/trasa/watchmud/command"
	"github.com/trasa/watchmud/event"
	"github.com/trasa/watchmud/gameserver"
)

func (w *World) handleInventory(msg *gameserver.HandlerParameter, cmd command.Inventory) {
	var items []event.InventoryItem
	for _, instPtr := range msg.Player.Inventory().GetAll() {
		if !msg.Player.Slots().IsItemInUse(instPtr) {
			items = append(items, event.InventoryItem{
				Id:               instPtr.Id.String(),
				ShortDescription: instPtr.Definition.ShortDescription,
				Categories:       instPtr.Definition.Categories.ToStringList(),
			})
		}
	}
	msg.Player.Send(event.Inventory{Items: items})
}
