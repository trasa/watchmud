package world

import (
	"github.com/trasa/watchmud/message"
)

func (w *World) handleSay(msg *message.IncomingMessage) {
	sayRequest := msg.Request.(message.SayRequest)
	room := w.playerRooms.playerToRoom[msg.Player]
	if room == nil {
		// player isn't in a room... not much to say really.
		msg.Player.Send(message.Response{
			MessageType: "say",
			Successful:  false,
			ResultCode:  "NOT_IN_A_ROOM",
		})
	} else {
		w.SendToPlayersInRoomExcept(msg.Player, room, message.SayNotification{
			Notification: message.Notification{MessageType: "say_notification"},
			Value:        sayRequest.Value,
			Sender:       msg.Player.GetName(),
		})
		msg.Player.Send(message.SayResponse{
			Response: message.Response{
				MessageType: "say",
				Successful:  true,
				ResultCode:  "OK",
			},
			Value: sayRequest.Value,
		})
	}
}
