package world

import (
	"github.com/trasa/watchmud/message"
)

// Tell everybody in the game something.
func (w *World) handleTellAll(msg *message.IncomingMessage) {
	tellAllRequest := msg.Request.(message.TellAllRequest)
	if tellAllRequest.Value != "" {
		w.SendToAllPlayersExcept(msg.Player, message.TellAllNotification{
			Response: message.NewSuccessfulResponse("tell_all_notification"),
			Value:    tellAllRequest.Value,
			Sender:   msg.Player.GetName(),
		})
		msg.Player.Send(message.NewSuccessfulResponse("tell_all"))
	} else {
		msg.Player.Send(message.NewUnsuccessfulResponse("tell_all", "NO_VALUE"))
	}
}
