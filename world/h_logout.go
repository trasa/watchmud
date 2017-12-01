package world

import (
	"github.com/trasa/watchmud/gameserver"
	"github.com/trasa/watchmud/message"
	"log"
)

func (w *World) handleLogout(msg *gameserver.HandlerParameter) {
	if msg.Player != nil {
		log.Printf("Player %s Logout", msg.Player.GetName())
		playerRoom := w.getRoomContainingPlayer(msg.Player)
		w.RemovePlayer(msg.Player)
		if playerRoom != nil {
			playerRoom.Send(message.LogoutNotification{
				Success:    true,
				ResultCode: "OK",
				PlayerName: msg.Player.GetName(),
			})
		}
	}
}
