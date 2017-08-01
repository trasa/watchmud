package world

import (
	"github.com/trasa/watchmud/direction"
	"github.com/trasa/watchmud/message"
	"log"
)

func (w *World) handleMove(msg *message.IncomingMessage) {
	// go somewhere
	playerRoom := w.getRoomContainingPlayer(msg.Player)
	// get the direction we want to go to
	dir := msg.Request.(message.MoveRequest).Direction

	// dirstr only used for log message so we'll ignore errors
	dirstr, _ := direction.DirectionToString(dir)
	log.Printf("player %s in room %s wants to move %s",
		msg.Player.GetName(),
		playerRoom.Name,
		dirstr,
	)

	// can player go in that direction?
	if targetRoom := playerRoom.Get(dir); targetRoom != nil {
		// make it happen
		w.movePlayer(msg.Player, dir, playerRoom, targetRoom)
		// send response message
		msg.Player.Send(message.MoveResponse{
			Response:        message.NewSuccessfulResponse("move"),
			RoomDescription: targetRoom.CreateRoomDescription(),
		})
	} else {
		// you can't go that way, tell player about error
		msg.Player.Send(message.Response{
			MessageType: "move",
			Successful:  false,
			ResultCode:  "CANT_GO_THAT_WAY",
		})
	}
}
