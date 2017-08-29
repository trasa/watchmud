package world

import (
	"github.com/trasa/watchmud/message"
)

func (w *World) handleTell(msg *message.IncomingMessage) {
	tellRequest := msg.Request.(message.TellRequest)
	sender := msg.Player.GetName()
	receiver := w.findPlayerByName(tellRequest.ReceiverPlayerName)
	value := tellRequest.Value

	if receiver == nil {
		msg.Player.Send(message.ResponseBase{
			MessageType: "tell",
			Successful:  false,
			ResultCode:  "TO_PLAYER_NOT_FOUND",
		})
	} else {
		receiver.Send(message.TellNotification{
			Response: message.ResponseBase{MessageType: "tell_notification"},
			Sender:       sender,
			Value:        value,
		})
		msg.Player.Send(message.ResponseBase{
			MessageType: "tell",
			Successful:  true,
			ResultCode:  "OK",
		})
	}
}
