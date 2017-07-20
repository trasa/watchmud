package world

import (
	"github.com/trasa/watchmud/message"
	"github.com/trasa/watchmud/response"
)

func (w *World) handleTell(message *message.IncomingMessage) {
	fromName := message.Player.GetName()
	toPlayer := w.findPlayerByName(message.Body["to"])
	value := message.Body["value"]

	if toPlayer == nil {
		message.Player.Send(response.Response{
			MessageType: "tell",
			Successful:  false,
			ResultCode:  "TO_PLAYER_NOT_FOUND",
		})
	} else {
		toPlayer.Send(response.TellNotification{
			Notification: response.Notification{MessageType: "tell"},
			From:         fromName,
			Value:        value,
		})
		message.Player.Send(response.Response{
			MessageType: "tell",
			Successful:  true,
			ResultCode:  "OK",
		})
	}
}
