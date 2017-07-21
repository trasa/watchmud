package world

import (
	"github.com/trasa/watchmud/message"
)

func (w *World) handleTell(msg *message.IncomingMessage) {
	tellRequest := msg.Request.(message.TellRequest)
	fromName := msg.Player.GetName()
	toPlayer := w.findPlayerByName(tellRequest.ReceiverPlayerName)
	value := tellRequest.Value

	if toPlayer == nil {
		msg.Player.Send(message.Response{
			MessageType: "tell",
			Successful:  false,
			ResultCode:  "TO_PLAYER_NOT_FOUND",
		})
	} else {
		toPlayer.Send(message.TellNotification{
			Notification: message.Notification{MessageType: "tell"},
			From:         fromName,
			Value:        value,
		})
		msg.Player.Send(message.Response{
			MessageType: "tell",
			Successful:  true,
			ResultCode:  "OK",
		})
	}
}
