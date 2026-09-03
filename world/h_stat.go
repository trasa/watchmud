package world

import (
	"log"

	message "github.com/trasa/watchmud-message"
	"github.com/trasa/watchmud/gameserver"
)

func (w *World) handleStat(msg *gameserver.HandlerParameter) {
	player := msg.Player
	if err := player.Send(message.StatResponse{
		Success:       true,
		ResultCode:    "OK",
		PlayerName:    player.GetName(),
		CurrentHealth: player.GetCurrentHealth(),
		MaxHealth:     player.GetMaxHealth(),
		Race:          "", // TODO rename race.RaceName,
		Class:         "", // TODO rename? class.ClassName,
		ZoneId:        player.Location().ZoneId,
		RoomId:        player.Location().RoomId,
		// TODO change to correct abilities
		Strength:     0,
		Dexterity:    0,
		Constitution: 0,
		Intelligence: 0,
		Wisdom:       0,
		Charisma:     0,
	}); err != nil {
		log.Printf("stat: Failed to send StatResponse to player %s: %v", player.GetName(), err)
	}
}
