package world

import (
	"github.com/trasa/watchmud/command"
	"github.com/trasa/watchmud/event"
	"github.com/trasa/watchmud/gameserver"
)

func (w *World) handleExits(msg *gameserver.HandlerParameter, cmd command.Exits) {
	r := w.getRoomContainingPlayer(msg.Player)
	if r == nil {
		r = w.VoidRoom
	}
	exits := []event.Exit{}
	for _, rexit := range r.Exits(false) {
		exits = append(exits, event.Exit{
			Direction: rexit.Direction,
			RoomName:  rexit.Room.Name,
		})
	}
	msg.Player.Send(event.Exits{Exits: exits})
}
