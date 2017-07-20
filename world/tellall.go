package world

import (
	"github.com/trasa/watchmud/message"
	"github.com/trasa/watchmud/response"
)

// Tell everybody in the game something.
func (w *World) handleTellAll(message *message.IncomingMessage) {
	if val, ok := message.Body["value"]; ok {
		w.SendToAllPlayersExcept(message.Player, response.TellAllNotification{
			Notification: response.Notification{
				MessageType: "tell_all_notification",
			},
			Value:  val,
			Sender: message.Player.GetName(),
		})
		message.Player.Send(response.Response{
			MessageType: "tell_all",
			Successful:  true,
			ResultCode:  "OK",
		})
	} else {
		message.Player.Send(response.Response{
			MessageType: "tell_all",
			Successful:  false,
			ResultCode:  "NO_VALUE",
		})
	}
}
