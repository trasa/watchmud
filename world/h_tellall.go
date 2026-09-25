package world

import (
	"github.com/watchmud/watchmud/command"
	"github.com/watchmud/watchmud/event"
	"github.com/watchmud/watchmud/gameserver"
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
