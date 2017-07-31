package world

import "github.com/trasa/watchmud/message"

func (w *World) handleLook(msg *message.IncomingMessage) {
	// for now, only "look" (no args) is supported
	// this will show the player the room they are in currently (if any)

	// get room for player
	playerRoom := w.GetRoomContainingPlayer(msg.Player)
	var resp message.LookResponse
	if playerRoom == nil {
		resp = w.VoidRoom.BuildLookResponse()
	} else {
		resp = playerRoom.BuildLookResponse()
	}
	msg.Player.Send(resp)
}
