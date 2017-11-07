package world

import (
	"github.com/trasa/watchmud/message"
)

func (w *World) handleExits(msg *message.IncomingMessage) {
	r := w.getRoomContainingPlayer(msg.Player)
	if r == nil {
		r = w.voidRoom
	}
	// convert directions to strings because json
	messageExitInfo := []message.DirectionToRoomName{}
	for _, rexit := range r.GetRoomExits(false) {
		messageExitInfo = append(messageExitInfo,
			message.DirectionToRoomName{
				Direction: rexit.Direction,
				RoomName:  rexit.Room.Name,
			})
	}
	resp := message.ExitsResponse{
		Response: message.NewSuccessfulResponse("exits"),
		ExitInfo: messageExitInfo,
	}
	msg.Player.Send(resp)
}
