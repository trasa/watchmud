package world

import (
	"github.com/trasa/watchmud-message"
	"github.com/trasa/watchmud/gameserver"
)

func (w *World) handleRecall(msg *gameserver.HandlerParameter) {

	// TODO determine if the player is allowed to do this command

	// TODO end combat?

	w.movePlayerMagically(msg.Player, w.StartRoom)

	// move the player to the "recall room"
	msg.Player.Send(message.RecallResponse{
		Success: true,
		ResultCode: "OK",
		RoomDescription: w.StartRoom.CreateRoomDescription(msg.Player),
	})
}
