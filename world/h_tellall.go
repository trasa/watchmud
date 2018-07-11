package world

import (
	"github.com/trasa/watchmud-message"
	"github.com/trasa/watchmud/gameserver"
)

// Tell everybody in the game something.
func (w *World) handleTellAll(msg *gameserver.HandlerParameter) {
	tellAllRequest := msg.Message.GetTellAllRequest()
	if tellAllRequest.Value != "" {
		w.SendToAllPlayersExcept(msg.Player, message.TellAllNotification{
			Success:    true,
			ResultCode: "OK",
			Value:      tellAllRequest.Value,
			Sender:     msg.Player.GetName(),
		})
		msg.Player.Send(message.TellAllResponse{Success: true, ResultCode: "OK"})
	} else {
		msg.Player.Send(message.TellAllResponse{Success: false, ResultCode: "NO_VALUE"})
	}
}
