package world

import (
	"github.com/trasa/watchmud/direction"
	"github.com/trasa/watchmud/message"
	"log"
)

func (w *World) handleGo(msg *message.IncomingMessage) {
	// go somewhere
	playerRoom := w.GetRoomContainingPlayer(msg.Player)
	// get the direction we want to go to
	dir := msg.Request.(message.GoRequest).Direction

	// dirstr only used for log message so we'll ignore errors
	dirstr, _ := direction.DirectionToString(dir)
	log.Printf("player %s in room %s wants to go %s",
		msg.Player.GetName(),
		playerRoom.Name,
		dirstr,
	)

	// can player go in that direction?
	if playerRoom.HasExit(dir) {
		// make it happen
	} else {
		// you can't go that way, tell player about error
	}
}
