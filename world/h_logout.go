package world

import (
	"github.com/trasa/watchmud/message"
	"log"
)

func (w *World) handleLogout(msg *message.IncomingMessage) {
	if msg.Player != nil {
		log.Printf("Player %s Logout", msg.Player.GetName())
		w.RemovePlayer(msg.Player)
	}
}
