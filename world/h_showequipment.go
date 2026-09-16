package world

import (
	"github.com/trasa/watchmud/command"
	"github.com/trasa/watchmud/event"
	"github.com/trasa/watchmud/gameserver"
)

func (w *World) handleShowEquipment(msg *gameserver.HandlerParameter, cmd command.ShowEquipment) {
	var items []event.EquippedItem
	// Equipment.All yields in slot order: the listing is the same every time.
	for loc, inst := range msg.Player.Equipment().All() {
		items = append(items, event.EquippedItem{
			Id:               inst.Id.String(),
			ShortDescription: inst.Definition.ShortDescription,
			Slot:             loc,
		})
	}
	msg.Player.Send(event.Equipment{Items: items})
}
