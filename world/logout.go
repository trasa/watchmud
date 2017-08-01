package world

import (
	"github.com/trasa/watchmud/message"
	"log"
)

func (w *World) handleLogout(msg *message.IncomingMessage) {
	// TODO not sure what to do here
	if msg.Player != nil {
		log.Printf("Player %s Logout", msg.Player.GetName())
	}
}
