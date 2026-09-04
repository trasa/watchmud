package world

import (
	"github.com/rs/zerolog/log"
	"github.com/trasa/watchmud/player"
)

func logWizCommand(p *player.Player, command string, msg string, args ...interface{}) {
	log.Warn().
		Str("playerId", p.Id).
		Str("playerName", p.Name).
		Str("commandType", "wiz").
		Str("command", command).
		Msgf(msg, args...)
}
