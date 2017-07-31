package world

import (
	"github.com/trasa/watchmud/direction"
	"github.com/trasa/watchmud/message"
	"log"
)

func (w *World) handleGo(msg *message.IncomingMessage) {
	// go somewhere
	playerRoom := w.GetRoomContainingPlayer(msg.Player)
	// get the direction we want to go to
	dir := msg.Request.(message.GoRequest).Direction

	// dirstr only used for log message so we'll ignore errors
	dirstr, _ := direction.DirectionToString(dir)
	log.Printf("player %s in room %s wants to go %s",
		msg.Player.GetName(),
		playerRoom.Name,
		dirstr,
	)

	// can player go in that direction?
	if targetRoom := playerRoom.Get(dir); targetRoom != nil {
		// make it happen
		// remove player from playerRoom (and tell everybody in playerRoom about it)
		w.LeaveRoom(msg.Player, playerRoom)
		// add player to playerRoom.direction()  (and tell everybody in that room about it)
		w.EnterRoom(msg.Player, targetRoom)
		// send response message
		msg.Player.Send(message.GoResponse{
			Response: message.Response{Successful: true,
				MessageType: "go",
				ResultCode:  "OK"},
			RoomDescription: targetRoom.BuildRoomDescription(),
		})
	} else {
		// you can't go that way, tell player about error
		msg.Player.Send(message.Response{
			MessageType: "go",
			Successful:  false,
			ResultCode:  "CANT_GO_THAT_WAY",
		})
	}
}
