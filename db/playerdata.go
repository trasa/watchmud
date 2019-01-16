package db

import (
	"github.com/trasa/watchmud/player"
	"log"
)

type PlayerData struct {
	Name      string `db:"player_name"`
	Inventory []PlayerInventoryData
	//slots     *object.Slots // maps instance ids to location
	CurHealth int64 `db:"current_health"`
	MaxHealth int64 `db:"max_health"`
	// current location of the player? "zoneId.roomId" ?
	Race  int32
	Class int32
}

const NewPlayerMaxHealth = 100

func GetPlayerData(playerName string) (result *PlayerData, err error) {
	log.Printf("DB - getting player data for %s", playerName)
	result = &PlayerData{}
	err = watchdb.Get(result, "SELECT player_name, current_health, max_health, race, class FROM players where player_name = $1", playerName)
	if err != nil {
		return
	}
	result.Inventory, err = getPlayerInventoryData(playerName)
	if err != nil {
		log.Printf("err %v", err)
	}

	for _, inv := range result.Inventory {
		log.Printf("inv %v", inv)
	}
	return
}

// Write the player back to the database
func SavePlayer(player player.Player) (err error) {
	log.Printf("DB - Save, examine dirty flag: %v", player.IsDirty())
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
	_, err = tx.NamedExec("UPDATE players SET current_health = :curHealth, max_health = :maxHealth where player_name = :name",
		map[string]interface{}{
			"name":      player.GetName(),
			"curHealth": player.GetCurrentHealth(),
			"maxHealth": player.GetMaxHealth(),
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

	// TODO save other stuff ...

	// success
	err = tx.Commit()
	return
}

func CreatePlayerData(playerName string, race int32, class int32) (result *PlayerData, err error) {
	log.Printf("DB - Create Player for %s", playerName)
	_, err = watchdb.NamedExec("INSERT INTO players (player_name, current_health, max_health, race, class) VALUES (:name, :curHealth, :maxHealth, :race, :class)",
		map[string]interface{}{
			"name":      playerName,
			"curHealth": NewPlayerMaxHealth,
			"maxHealth": NewPlayerMaxHealth,
			"race":      race,
			"class":     class,
		})

	if err == nil {
		result = &PlayerData{
			Name:      playerName,
			CurHealth: NewPlayerMaxHealth,
			MaxHealth: NewPlayerMaxHealth,
			Race:      race,
			Class:     class,
		}
	}
	// TODO new player equipment, inventory ...
	return
}
