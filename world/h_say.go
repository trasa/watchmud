package world

import (
	"github.com/watchmud/watchmud/command"
	"github.com/watchmud/watchmud/event"
	"github.com/watchmud/watchmud/gameserver"
)

func (w *World) handleSay(msg *gameserver.HandlerParameter, cmd command.Say) {
	room := w.playerRoom(msg.Player)
	room.Send(event.Said{
		Speaker: msg.Player.Name(),
		Value:   cmd.Value,
	})
}
