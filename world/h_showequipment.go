package world

import (
	"github.com/watchmud/watchmud/command"
	"github.com/watchmud/watchmud/event"
	"github.com/watchmud/watchmud/gameserver"
)

func (w *World) handleShowEquipment(msg *gameserver.HandlerParameter, cmd command.ShowEquipment) {
	var items []event.EquippedItem
	// Equipment.All yields in slot order: the listing is the same every time.
	for loc, inst := range msg.Player.Equipment().All() {
		items = append(items, event.EquippedItem{
			Id:               inst.Id.String(),
			ShortDescription: inst.Definition.ShortDescription,
			Slot:             loc,
			Durability:       inst.Durability,
			MaxDurability:    inst.Definition.MaxDurability,
			Broken:           inst.Broken(),
			Power:            inst.Power,
		})
	}
	msg.Player.Send(event.Equipment{Power: msg.Player.Power(), Items: items})
}
