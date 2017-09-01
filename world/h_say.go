package world

import (
	"github.com/trasa/watchmud/message"
)

func (w *World) handleSay(msg *message.IncomingMessage) {
	sayRequest := msg.Request.(message.SayRequest)
	room := w.playerRooms.playerToRoom[msg.Player]
	if room == nil {
		// player isn't in a room... not much to say really.
		msg.Player.Send(message.ResponseBase{
			MessageType: "say",
			Successful:  false,
			ResultCode:  "NOT_IN_A_ROOM",
		})
	} else {
		room.SendExcept(msg.Player, message.SayNotification{
			Response: message.NewSuccessfulResponse("say_notification"),
			Value:    sayRequest.Value,
			Sender:   msg.Player.GetName(),
		})
		msg.Player.Send(message.SayResponse{
			Response: message.NewSuccessfulResponse("say"),
			Value:    sayRequest.Value,
		})
	}
}
