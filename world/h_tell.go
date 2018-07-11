package world

import (
	"github.com/trasa/watchmud-message"
	"github.com/trasa/watchmud/gameserver"
)

func (w *World) handleTell(msg *gameserver.HandlerParameter) {
	tellRequest := msg.Message.GetTellRequest()
	sender := msg.Player.GetName()
	receiver := w.findPlayerByName(tellRequest.ReceiverPlayerName)
	value := tellRequest.Value

	if receiver == nil {
		msg.Player.Send(message.TellResponse{Success: false, ResultCode: "TO_PLAYER_NOT_FOUND"})
	} else {
		receiver.Send(message.TellNotification{
			Success:    true,
			ResultCode: "OK",
			Sender:     sender,
			Value:      value,
		})
		msg.Player.Send(message.TellResponse{Success: true, ResultCode: "OK"})
	}
}
