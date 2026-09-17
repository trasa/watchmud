package world

import (
	"github.com/trasa/watchmud/command"
	"github.com/trasa/watchmud/event"
	"github.com/trasa/watchmud/gameserver"
)

func (w *World) handleStat(msg *gameserver.HandlerParameter, cmd command.Stat) {
	p := msg.Player

	// The room the world has the player in, not p.Location(): the location on
	// the player is only ever set at load time and movePlayer doesn't update
	// it. See ROADMAP.md "Dual location bookkeeping".

	room := w.getPlayerRoom(p)

	p.Send(event.Stat{
		PlayerName:    p.Name(),
		Lineage:       p.LineageName(),
		Role:          w.roleName(p.RoleWeights()),
		CurrentHealth: p.CurrentHealth(),
		MaxHealth:     p.MaxHealth(),
		ZoneId:        room.Zone.Id,
		RoomId:        room.Id,
	})
}
