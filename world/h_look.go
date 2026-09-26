package world

import (
	"github.com/watchmud/watchmud/command"
	"github.com/watchmud/watchmud/gameserver"
)

func (w *World) handleLook(msg *gameserver.HandlerParameter, cmd command.Look) {
	if cmd.In {
		w.lookIn(msg, cmd.Target)
		return
	}
	// for now, only "look" (no args) is supported
	// this will show the player the room they are in currently (if any)
	playerRoom := w.playerRoom(msg.Player)
	msg.Player.Send(playerRoom.DescriptionExcept(msg.Player))
}
