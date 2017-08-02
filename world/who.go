package world

import (
	"github.com/trasa/watchmud/message"
	"github.com/trasa/watchmud/player"
)

func (w *World) handleWho(msg *message.IncomingMessage) {
	// in the future we'll need to split this up by
	// rank, security, other things, but for now show
	// everybody everything.

	// playerName, (level, class, other things we don't have yet), zoneName, roomName
	info := []message.WhoPlayerInfo{}
	w.PlayerList.Iter(func(p player.Player) {
		r := w.getRoomContainingPlayer(p)
		var zoneName, roomName string
		if r == nil {
			zoneName = r.Zone.Name
			roomName = r.Name
		} else {
			zoneName = "(None)"
			roomName = "(None)"
		}
		info = append(info, message.WhoPlayerInfo{
			PlayerName: p.GetName(),
			ZoneName:   zoneName,
			RoomName:   roomName,
		})
	})

	msg.Player.Send(message.WhoResponse{
		Response:   message.NewSuccessfulResponse("who"),
		PlayerInfo: info,
	})
}
