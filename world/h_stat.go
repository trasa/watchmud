package world

import (
	"github.com/trasa/watchmud/command"
	"github.com/trasa/watchmud/event"
	"github.com/trasa/watchmud/gameserver"
)

func (w *World) handleStat(msg *gameserver.HandlerParameter, cmd command.Stat) {
	player := msg.Player
	player.Send(event.Stat{
		PlayerName:    player.Name(),
		CurrentHealth: 0,  // TODO
		MaxHealth:     0,  // TODO
		Lineage:       "", // TODO Phase 6: rules.Catalog lookup
		Class:         "", // TODO Phase 6: class.ClassName
		ZoneId:        "", // TODO player.Location.ZoneId
		RoomId:        "", // TODO player.Location.RoomId
		// TODO change to correct abilities
		Strength:     0,
		Dexterity:    0,
		Constitution: 0,
		Intelligence: 0,
		Wisdom:       0,
		Charisma:     0,
	})
}
