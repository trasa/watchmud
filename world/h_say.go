package world

import (
	"github.com/trasa/watchmud/command"
	"github.com/trasa/watchmud/event"
	"github.com/trasa/watchmud/gameserver"
)

func (w *World) handleSay(msg *gameserver.HandlerParameter, cmd command.Say) {
	room := w.playerRooms.playerToRoom[msg.Player]
	if room == nil {
		// player isn't in a room... not much to say really.
		msg.Fail(event.NotInARoom)
		return
	}
	room.Send(event.Said{
		Speaker: msg.Player.Name(),
		Value:   cmd.Value,
	})
}
