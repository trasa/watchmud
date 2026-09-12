package world

import (
	"github.com/trasa/watchmud/command"
	"github.com/trasa/watchmud/event"
	"github.com/trasa/watchmud/gameserver"
)

// Tell everybody in the game something.
func (w *World) handleTellAll(msg *gameserver.HandlerParameter, cmd command.TellAll) {
	if cmd.Value == "" {
		msg.Fail(event.NoValue)
		return
	}
	shouted := event.Shouted{
		Speaker: msg.Player.Name(),
		Value:   cmd.Value,
	}
	w.SendToAllPlayersExcept(msg.Player, shouted)
	msg.Player.Send(shouted)
}
