package world

import (
	"github.com/trasa/watchmud-message"
	"github.com/trasa/watchmud-message/direction"
	"github.com/trasa/watchmud/gameserver"
	"log"
)

func (w *World) handleMove(msg *gameserver.HandlerParameter) {

	// is player in a fight?
	if w.fightLedger.IsFighting(msg.Player) || w.fightLedger.IsBeingFought(msg.Player) {
		msg.Player.Send(message.MoveResponse{
			Success:    false,
			ResultCode: "IN_A_FIGHT",
		})
		return
	}

	// go somewhere
	playerRoom := w.getRoomContainingPlayer(msg.Player)
	// get the direction we want to go to
	dir := direction.Direction(msg.Message.GetMoveRequest().Direction)

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
			Success:         true,
			ResultCode:      "OK",
			RoomDescription: targetRoom.CreateRoomDescription(msg.Player),
		})
	} else {
		// you can't go that way, tell player about error
		msg.Player.Send(message.MoveResponse{
			Success:    false,
			ResultCode: "CANT_GO_THAT_WAY",
		})
	}
}
