package world

import (
	"sort"

	"github.com/trasa/watchmud/command"
	"github.com/trasa/watchmud/event"
	"github.com/trasa/watchmud/gameserver"
	"github.com/trasa/watchmud/player"
)

func (w *World) handleWho(msg *gameserver.HandlerParameter, cmd command.Who) {
	// in the future we'll need to split this up by
	// rank, security, other things, but for now show
	// everybody everything.

	// playerName, (level, class, other things we don't have yet), zoneName, roomName
	entries := []event.WhoEntry{}
	w.playerList.Iter(func(p *player.Player) {
		r := w.getRoomContainingPlayer(p)
		var zoneName, roomName string
		if r != nil {
			zoneName = r.Zone.Name
			roomName = r.Name
		}
		entries = append(entries, event.WhoEntry{
			PlayerName: p.Name(),
			ZoneName:   zoneName,
			RoomName:   roomName,
		})
	})

	// sort results by name
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].PlayerName < entries[j].PlayerName
	})

	msg.Player.Send(event.Who{Players: entries})
}
