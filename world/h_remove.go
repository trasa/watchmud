package world

import (
	"github.com/trasa/watchmud/command"
	"github.com/trasa/watchmud/event"
	"github.com/trasa/watchmud/gameserver"
)

// handleRemove takes an equipped item off. The item stays in the player's
// inventory -- wearing something never took it out of there in the first
// place -- so this only empties the slot.
//
// Emptying the slot is what makes a role switchable: the weights it was
// contributing stop counting the moment this returns.
func (w *World) handleRemove(msg *gameserver.HandlerParameter, cmd command.Remove) {
	if cmd.Target == "" {
		msg.Fail(event.NoTarget)
		return
	}

	loc, inst, found := msg.Player.Slots().FindByNameOrAlias(cmd.Target)
	if !found {
		msg.Fail(event.TargetNotFound)
		return
	}

	msg.Player.Slots().Clear(loc)
	msg.Player.Send(event.Removed{Item: inst.Definition.ShortDescription})
}
