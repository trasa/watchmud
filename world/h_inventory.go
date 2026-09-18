package world

import (
	"github.com/trasa/watchmud/command"
	"github.com/trasa/watchmud/event"
	"github.com/trasa/watchmud/gameserver"
)

func (w *World) handleInventory(msg *gameserver.HandlerParameter, cmd command.Inventory) {
	var items []event.InventoryItem
	for instPtr := range msg.Player.Inventory().All() {
		if !msg.Player.Equipment().ItemEquipped(instPtr) {
			items = append(items, event.InventoryItem{
				Id:               instPtr.Id.String(),
				ShortDescription: instPtr.Definition.ShortDescription,
				Category:         instPtr.Definition.ObjectCategory,
			})
		}
	}
	msg.Player.Send(event.Inventory{Items: items})
}
