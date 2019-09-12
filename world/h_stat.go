package world

import (
	message "github.com/trasa/watchmud-message"
	"github.com/trasa/watchmud/db"
	"github.com/trasa/watchmud/gameserver"
	"log"
)

func (w *World) handleStat(msg *gameserver.HandlerParameter) {
	player := msg.Player
	// TODO cache these? don't want to hit the database
	race, err := db.GetSingleRaceData(player.GetRaceId())
	if err != nil {
		log.Printf("stat: Failed to retrieve race for player %s: %v", player.GetName(), err)
		_ = player.Send(message.StatResponse{
			Success:    false,
			ResultCode: "RACE_DB_ERROR",
		})
	}
	class, err := db.GetSingleClassData(player.GetClassId())
	if err != nil {
		log.Printf("stat: Failed to retrieve class for player %s: %v", player.GetName(), err)
		_ = player.Send(message.StatResponse{
			Success:    false,
			ResultCode: "CLASS_DB_ERROR",
		})
	}

	if err := player.Send(message.StatResponse{
		Success:       true,
		ResultCode:    "OK",
		PlayerName:    player.GetName(),
		CurrentHealth: player.GetCurrentHealth(),
		MaxHealth:     player.GetMaxHealth(),
		Race:          race.RaceName,
		Class:         class.ClassName,
		ZoneId:        player.Location().ZoneId,
		RoomId:        player.Location().RoomId,
		Strength:      int32(player.Abilities().Strength),
		Dexterity:     int32(player.Abilities().Dexterity),
		Constitution:  int32(player.Abilities().Constitution),
		Intelligence:  int32(player.Abilities().Intelligence),
		Wisdom:        int32(player.Abilities().Wisdom),
		Charisma:      int32(player.Abilities().Charisma),
	}); err != nil {
		log.Printf("stat: Failed to send StatResponse to player %s: %v", player.GetName(), err)
	}
}
