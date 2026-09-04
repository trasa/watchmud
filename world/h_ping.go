package world

import (
	"github.com/rs/zerolog/log"
	"github.com/trasa/watchmud-message"
	"github.com/trasa/watchmud/gameserver"
)

func (w *World) handlePing(msg *gameserver.HandlerParameter) {
	// TODO remove this constant, use logging parameters intead
	if VERBOSE_LOGGING {
		log.Trace().Msgf("Player %s Ping", msg.Player.Name)
	}
	// TODO error handling
	msg.Player.Send(message.Pong{
		Success:    true,
		ResultCode: "OK",
		Target:     msg.Message.GetPing().Target,
	})
}
