package world

import (
	"github.com/trasa/watchmud/command"
	"github.com/trasa/watchmud/gameserver"
)

func (w *World) handleRecall(msg *gameserver.HandlerParameter, cmd command.Recall) {

	// TODO determine if the player is allowed to do this command

	// TODO end combat?

	w.movePlayerMagically(msg.Player, w.StartRoom)

	// move the player to the "recall room"
	msg.Player.Send(w.StartRoom.DescriptionExcept(msg.Player))
}
