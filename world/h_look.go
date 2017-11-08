package world

import (
	"github.com/trasa/watchmud/gameserver"
	"github.com/trasa/watchmud/message"
)

func (w *World) handleLook(msg *gameserver.HandlerParameter) {
	// for now, only "look" (no args) is supported
	// this will show the player the room they are in currently (if any)

	// get room for player
	playerRoom := w.getRoomContainingPlayer(msg.Player)
	resp := message.LookResponse{
		Success: true, ResultCode: "OK",
	}
	if playerRoom == nil {
		playerRoom = w.voidRoom
	}
	resp.RoomDescription = playerRoom.CreateRoomDescription(msg.Player)
	msg.Player.Send(resp)
}
