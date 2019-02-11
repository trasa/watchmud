package db

import (
	"github.com/trasa/watchmud/player"
	"log"
)

type PlayerData struct {
	Id        int64  `db:"player_id"`
	Name      string `db:"player_name"`
	Inventory []PlayerInventoryData
	CurHealth int64 `db:"current_health"`
	MaxHealth int64 `db:"max_health"`
	// current location of the player? "zoneId.roomId" ?
	Race       int32        `db:"race"`
	Class      int32        `db:"class"`
	Slots      SlotDataList `db:"slots" json:"Slots"`
	LastZoneId *string      `db:"last_zone_id"`
	LastRoomId *string      `db:"last_room_id"`
}

const NewPlayerMaxHealth = 100

func GetPlayerData(playerName string) (result PlayerData, err error) {
	log.Printf("DB - getting player data for %s", playerName)
	result = PlayerData{}
	err = watchdb.Get(&result, "SELECT player_id, player_name, current_health, max_health, race, class, last_zone_id, last_room_id, slots FROM players where player_name = $1", playerName)
	if err != nil {
		return
	}

	result.Inventory, err = getPlayerInventoryData(result.Id)
	if err != nil {
		log.Printf("err %v", err)
	}

	log.Printf("Player %s slots: %v", playerName, len(result.Slots.Slots))
	return
}

// Write the player back to the database
func SavePlayer(player player.Player) (err error) {
	//log.Printf("DB - Save, examine dirty flag: %v", player.IsDirty())
	if !player.IsDirty() {
		return
	}

	log.Printf("DB - Saving Player Data for %s", player.GetName())
	tx, err := watchdb.Beginx()
	if err != nil {
		log.Printf("Error starting transaction to save player: %s - %s", player, err)
		if tx != nil {
			tx.Rollback()
		}
		return
	}

	// players table
	log.Printf("loc %s %s", player.Location().ZoneId, player.Location().RoomId)
	// TODO should only be updating zone and room location "every so often"
	// TODO also, updates should be queued up and executed in a different goroutine
	_, err = tx.NamedExec("UPDATE players SET "+
		"current_health = :curHealth, "+
		"max_health = :maxHealth, "+
		"last_zone_id = :lastZoneId, "+
		"last_room_id = :lastRoomId, "+
		"slots = :slots "+
		"where player_id = :id",
		map[string]interface{}{
			"id":         player.GetId(),
			"curHealth":  player.GetCurrentHealth(),
			"maxHealth":  player.GetMaxHealth(),
			"lastZoneId": player.Location().ZoneId,
			"lastRoomId": player.Location().RoomId,
			"slots":      NewSlotDataList(player.Slots()),
		})
	if err != nil {
		log.Printf("Failed to update players for player %s - %s", player, err)
		tx.Rollback()
		return
	}

	// player inventory
	err = savePlayerInventory(tx, player)
	if err != nil {
		log.Printf("Error saving player inventory %s - %s", player, err)
		tx.Rollback()
		return
	}

	// success
	err = tx.Commit()
	return
}

func CreatePlayerData(playerName string, race int32, class int32, zoneId string, roomId string) (result *PlayerData, err error) {
	log.Printf("DB - Create Player for %s", playerName)
	res, err := watchdb.NamedExec("INSERT INTO players (player_name, current_health, max_health, race, class, last_zone_id, last_room_id) "+
		"VALUES (:name, :curHealth, :maxHealth, :race, :class, :lastZoneId, :lastRoomId)",
		map[string]interface{}{
			"name":       playerName,
			"curHealth":  NewPlayerMaxHealth,
			"maxHealth":  NewPlayerMaxHealth,
			"race":       race,
			"class":      class,
			"lastZoneId": zoneId,
			"lastRoomId": roomId,
		})
	result.Id, err = res.LastInsertId()

	if err == nil {
		result = &PlayerData{
			Name:       playerName,
			CurHealth:  NewPlayerMaxHealth,
			MaxHealth:  NewPlayerMaxHealth,
			Race:       race,
			Class:      class,
			LastZoneId: &zoneId,
			LastRoomId: &roomId,
		}
	}
	// TODO new player equipment, inventory ...
	return
}
