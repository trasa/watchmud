package world

import "github.com/trasa/watchmud/message"

func (w *World) handleGo(msg *message.IncomingMessage) {
	// go somewhere
	playerRoom := w.GetRoomContainingPlayer(msg.Player)
	// get the direction we want to go to
	direction := msg.Request.(message.GoRequest).Direction

	// can player go in that direction?
	if playerRoom.HasExit(direction) {
		// make it happen
	} else {
		// you can't go that way, tell player about error
	}
}
