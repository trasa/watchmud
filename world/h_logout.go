package world

import (
	"log"

	"github.com/trasa/watchmud-message"
	"github.com/trasa/watchmud/gameserver"
)

func (w *World) handleLogout(msg *gameserver.HandlerParameter) {
	if msg.Player == nil {
		return
	}
	log.Printf("Player %s Logout", msg.Player.Name)
	playerRoom := w.getRoomContainingPlayer(msg.Player)
	w.RemovePlayer(msg.Player)
	if playerRoom != nil {
		playerRoom.Send(message.LogoutNotification{
			Success:    true,
			ResultCode: "OK",
			PlayerName: msg.Player.Name,
		})
	}
	// TODO reimplement saving
	/*
		log.Printf("final loc %s - %s", msg.Player.Location().ZoneId, msg.Player.Location().RoomId)
		if err := db.ForceSavePlayer(msg.Player); err != nil {
			log.Printf("Error saving player %s on logout - %s", msg.Player.GetName(), err)
		}*/
}
