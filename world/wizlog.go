package world

import (
	"github.com/rs/zerolog/log"
	"github.com/trasa/watchmud/player"
)

func logWizCommand(p player.Player, command string, msg string, args ...interface{}) {
	log.Warn().
		Int64("playerId", p.GetId()).
		Str("playerName", p.GetName()).
		Str("commandType", "wiz").
		Str("command", command).
		Msgf(msg, args...)
}
