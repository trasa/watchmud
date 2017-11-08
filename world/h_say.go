package world

import (
	"github.com/trasa/watchmud/gameserver"
	"github.com/trasa/watchmud/message"
)

func (w *World) handleSay(msg *gameserver.HandlerParameter) {
	sayRequest := msg.Message.GetSayRequest()
	room := w.playerRooms.playerToRoom[msg.Player]
	if room == nil {
		// player isn't in a room... not much to say really.
		msg.Player.Send(message.SayResponse{
			Success:    false,
			ResultCode: "NOT_IN_A_ROOM",
		})
	} else {
		room.SendExcept(msg.Player, message.SayNotification{
			Success:    true,
			ResultCode: "OK",
			Value:      sayRequest.Value,
			Sender:     msg.Player.GetName(),
		})
		msg.Player.Send(message.SayResponse{
			Success:    true,
			ResultCode: "OK",
			Value:      sayRequest.Value,
		})
	}
}
