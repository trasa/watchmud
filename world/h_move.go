package world

import (
	"github.com/rs/zerolog/log"
	"github.com/trasa/watchmud-message"
	"github.com/trasa/watchmud-message/direction"
	"github.com/trasa/watchmud/gameserver"
)

func (w *World) handleMove(msg *gameserver.HandlerParameter) {

	// is player in a fight?
	// TODO reimplement
	/*
		if w.fightLedger.IsFighting(msg.Player) || w.fightLedger.IsBeingFought(msg.Player) {
			msg.Player.Send(message.MoveResponse{
				Success:    false,
				ResultCode: "IN_A_FIGHT",
			})
			return
		}
	*/
	playerRoom := w.getRoomContainingPlayer(msg.Player)
	dir := direction.Direction(msg.Message.GetMoveRequest().Direction)

	log.Trace().Str("player", msg.Player.Name).Str("room", playerRoom.Name).Msgf("player wants to move %s", dir.String())

	// can player go in that direction?
	if targetRoom := playerRoom.Get(dir); targetRoom != nil {
		// make it happen
		w.movePlayer(msg.Player, dir, playerRoom, targetRoom)
		msg.Player.Send(message.MoveResponse{
			Success:         true,
			ResultCode:      "OK",
			RoomDescription: targetRoom.CreateRoomDescription(msg.Player),
		})
	} else {
		msg.Player.Send(message.MoveResponse{
			Success:    false,
			ResultCode: "CANT_GO_THAT_WAY",
		})
	}
}
