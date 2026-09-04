package world

import (
	"github.com/rs/zerolog/log"
	message "github.com/trasa/watchmud-message"
	"github.com/trasa/watchmud/gameserver"
)

func (w *World) handleStat(msg *gameserver.HandlerParameter) {
	player := msg.Player
	if err := player.Send(message.StatResponse{
		Success:       true,
		ResultCode:    "OK",
		PlayerName:    player.Name,
		CurrentHealth: 0,  // TODO
		MaxHealth:     0,  // TODO
		Race:          "", // TODO rename Lineage
		Class:         "", // TODO rename? class.ClassName,
		ZoneId:        "", // TODO player.Location.ZoneId,
		RoomId:        "", // TODO player.Location.RoomId,
		// TODO change to correct abilities
		Strength:     0,
		Dexterity:    0,
		Constitution: 0,
		Intelligence: 0,
		Wisdom:       0,
		Charisma:     0,
	}); err != nil {
		log.Error().Msgf("stat: Failed to send StatResponse to player %s: %v", player.Name, err)
	}
}
