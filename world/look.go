package world

import "github.com/trasa/watchmud/message"

func (w *World) handleLook(msg *message.IncomingMessage) {
	//lookRequest := msg.Request.(message.LookRequest)
	// for now, only "look" (no args) is supported
	// this will show the player the room they are in currently (if any)

	// get room for player
	playerRoom := w.GetRoomContainingPlayer(msg.Player)
	var resp message.LookResponse
	if playerRoom == nil {
		resp = message.LookResponse{
			Response: message.Response{MessageType: "look", Successful: true},
			RoomName: "Not In A Room",
			Value:    "You see nothing but endless void.",
		}
	} else {
		resp = message.LookResponse{
			Response: message.Response{MessageType: "look", Successful: true},
			RoomName: playerRoom.Name,
			Value:    playerRoom.Description,
			Exits:    playerRoom.GetExits(),
			// TODO other occupants or objects in the room
		}
	}
	msg.Player.Send(resp)
}
