package world

import (
	"github.com/trasa/watchmud/command"
	"github.com/trasa/watchmud/event"
	"github.com/trasa/watchmud/gameserver"
)

func (w *World) handleSay(msg *gameserver.HandlerParameter, cmd command.Say) {
	room := w.getPlayerRoom(msg.Player)
	room.Send(event.Said{
		Speaker: msg.Player.Name(),
		Value:   cmd.Value,
	})
}
