package world

import (
	"github.com/trasa/watchmud-message"
	"github.com/trasa/watchmud/gameserver"
	"github.com/trasa/watchmud/player"
	"sort"
)

func (w *World) handleWho(msg *gameserver.HandlerParameter) {
	// in the future we'll need to split this up by
	// rank, security, other things, but for now show
	// everybody everything.

	// playerName, (level, class, other things we don't have yet), zoneName, roomName
	info := []*message.WhoResponse_PlayerInfo{}
	w.playerList.Iter(func(p player.Player) {
		r := w.getRoomContainingPlayer(p)
		var zoneName, roomName string
		if r != nil {
			zoneName = r.Zone.Name
			roomName = r.Name
		}
		info = append(info, &message.WhoResponse_PlayerInfo{
			PlayerName: p.GetName(),
			ZoneName:   zoneName,
			RoomName:   roomName,
		})
	})

	// sort results by name
	sort.Slice(info, func(i, j int) bool {
		return info[i].PlayerName < info[j].PlayerName
	})

	msg.Player.Send(message.WhoResponse{
		Success:    true,
		ResultCode: "OK",
		PlayerInfo: info,
	})
}
