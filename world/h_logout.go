package world

import (
	"github.com/trasa/watchmud/command"
	"github.com/trasa/watchmud/event"
	"github.com/trasa/watchmud/gameserver"
)

func (w *World) handleLogout(msg *gameserver.HandlerParameter, cmd command.Logout) {
	if msg.Player == nil {
		return
	}
	msg.Player.Log().Debug().Str("cause", cmd.Cause).Msg("player logout")
	// get the room before removing from the world
	room := w.playerToRoom.Get(msg.Player)
	if room == nil {
		return // already logged out or not in the world
	}
	w.RemovePlayer(msg.Player)

	// the player is already out of the room, so this reaches everyone else
	room.Send(event.LoggedOut{Actor: msg.Player.Name()})

	if err := w.store.Save(msg.Player.Record()); err != nil {
		msg.Player.Log().Err(err).Msg("Error saving player on logout")
	}
}
