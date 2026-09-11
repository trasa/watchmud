package world

import (
	"github.com/rs/zerolog/log"
	"github.com/trasa/watchmud-message"
	"github.com/trasa/watchmud/gameserver"
)

func (w *World) handleLogout(msg *gameserver.HandlerParameter) /*error */ {
	// TODO: need to add error handling
	if msg.Player == nil {
		return /*nil*/
	}
	log.Info().Msgf("Player %s Logout", msg.Player.Name())
	playerRoom := w.getRoomContainingPlayer(msg.Player)
	w.RemovePlayer(msg.Player)
	if playerRoom != nil {
		playerRoom.Send(message.LogoutNotification{
			Success:    true,
			ResultCode: "OK",
			PlayerName: msg.Player.Name(),
		})
	}
	if err := w.store.Save(msg.Player.Record()); err != nil {
		log.Error().Err(err).Msg("Error saving player on logout")
		return /*err*/
	}
	return /*nil*/
}
