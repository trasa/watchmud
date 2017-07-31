package world

import "github.com/trasa/watchmud/message"

func (w *World) handleLook(msg *message.IncomingMessage) {
	// for now, only "look" (no args) is supported
	// this will show the player the room they are in currently (if any)

	// get room for player
	playerRoom := w.GetRoomContainingPlayer(msg.Player)
	resp := message.LookResponse{
		Response: message.NewSuccessfulResponse("look"),
	}
	if playerRoom == nil {
		resp.RoomDescription = w.VoidRoom.BuildRoomDescription()
	} else {
		resp.RoomDescription = playerRoom.BuildRoomDescription()
	}
	msg.Player.Send(resp)
}
