package world

import (
	"github.com/trasa/watchmud-message"
	"github.com/trasa/watchmud/db"
	"github.com/trasa/watchmud/gameserver"
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
		log.Printf("final loc %s - %s", msg.Player.Location().ZoneId, msg.Player.Location().RoomId)
		if err := db.ForceSavePlayer(msg.Player); err != nil {
			log.Printf("Error saving player %s on logout - %s", msg.Player.GetName(), err)
		}
	}
}
