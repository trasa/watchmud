package world

import (
	"github.com/trasa/watchmud/command"
	"github.com/trasa/watchmud/event"
	"github.com/trasa/watchmud/gameserver"
)

func (w *World) handleTell(msg *gameserver.HandlerParameter, cmd command.Tell) {
	receiver := w.findPlayerByName(cmd.To)
	if receiver == nil {
		msg.Fail(event.ToPlayerNotFound)
		return
	}
	// one event, both ends of the conversation
	told := event.Told{
		From:  msg.Player.Name(),
		To:    receiver.Name(),
		Value: cmd.Value,
	}
	receiver.Send(told)
	msg.Player.Send(told)
}
