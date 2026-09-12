package world

import (
	"github.com/rs/zerolog/log"
	"github.com/trasa/watchmud/command"
	"github.com/trasa/watchmud/event"
	"github.com/trasa/watchmud/gameserver"
)

func (w *World) handleLogout(msg *gameserver.HandlerParameter, cmd command.Logout) {
	if msg.Player == nil {
		return
	}
	log.Info().Str("player", msg.Player.Name()).Str("cause", cmd.Cause).Msg("player logout")
	playerRoom := w.getRoomContainingPlayer(msg.Player)
	w.RemovePlayer(msg.Player)
	if playerRoom != nil {
		// the player is already out of the room, so this reaches everyone else
		playerRoom.Send(event.LoggedOut{Actor: msg.Player.Name()})
	}
	if err := w.store.Save(msg.Player.Record()); err != nil {
		log.Error().Err(err).Msg("Error saving player on logout")
	}
}
