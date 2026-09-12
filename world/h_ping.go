package world

import (
	"github.com/rs/zerolog/log"
	"github.com/trasa/watchmud/command"
	"github.com/trasa/watchmud/event"
	"github.com/trasa/watchmud/gameserver"
)

func (w *World) handlePing(msg *gameserver.HandlerParameter, cmd command.Ping) {
	// TODO remove this constant, use logging parameters instead
	if VERBOSE_LOGGING {
		log.Trace().Msgf("Player %s Ping", msg.Player.Name())
	}
	msg.Player.Send(event.Pong{Target: cmd.Target})
}
