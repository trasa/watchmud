package world

import (
	"github.com/trasa/watchmud/gameserver"
	"log"
)

func (w *World) handleLogout(msg *gameserver.HandlerParameter) {
	if msg.Player != nil {
		log.Printf("Player %s Logout", msg.Player.GetName())
		w.RemovePlayer(msg.Player)
	}
}
