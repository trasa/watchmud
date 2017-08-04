package world

import "github.com/trasa/watchmud/message"

func (w *World) handleLook(msg *message.IncomingMessage) {
	// for now, only "look" (no args) is supported
	// this will show the player the room they are in currently (if any)

	// get room for player
	playerRoom := w.getRoomContainingPlayer(msg.Player)
	resp := message.LookResponse{
		Response: message.NewSuccessfulResponse("look"),
	}
	if playerRoom == nil {
		playerRoom = w.voidRoom
	}
	resp.RoomDescription = playerRoom.CreateRoomDescription()
	// add players to the room description
	// TODO only add those players you can actually see
	for _, p := range w.playerRooms.GetPlayers(playerRoom) {
		// don't add yourself
		if msg.Player != p {
			resp.RoomDescription.Players = append(resp.RoomDescription.Players, p.GetName())
		}
	}
	msg.Player.Send(resp)
}
