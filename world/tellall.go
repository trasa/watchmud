package world

import (
	"github.com/trasa/watchmud/message"
)

// Tell everybody in the game something.
func (w *World) handleTellAll(msg *message.IncomingMessage) {
	if val, ok := msg.Body["value"]; ok {
		w.SendToAllPlayersExcept(msg.Player, message.TellAllNotification{
			Notification: message.Notification{
				MessageType: "tell_all_notification",
			},
			Value:  val,
			Sender: msg.Player.GetName(),
		})
		msg.Player.Send(message.Response{
			MessageType: "tell_all",
			Successful:  true,
			ResultCode:  "OK",
		})
	} else {
		msg.Player.Send(message.Response{
			MessageType: "tell_all",
			Successful:  false,
			ResultCode:  "NO_VALUE",
		})
	}
}
