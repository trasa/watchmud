package world

import (
	"github.com/trasa/watchmud/command"
	"github.com/trasa/watchmud/gameserver"
)

func (w *World) handleLook(msg *gameserver.HandlerParameter, cmd command.Look) {
	// for now, only "look" (no args) is supported
	// this will show the player the room they are in currently (if any)
	playerRoom := w.getRoomContainingPlayer(msg.Player)
	if playerRoom == nil {
		playerRoom = w.VoidRoom
	}
	msg.Player.Send(playerRoom.DescriptionExcept(msg.Player))
}
