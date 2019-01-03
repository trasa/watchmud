package db

import (
	"github.com/jmoiron/sqlx"
	"github.com/satori/go.uuid"
	"github.com/trasa/watchmud/player"
	"log"
)

type PlayerInventoryData struct {
	PlayerName   string    `db:"player_name"`
	InstanceId   uuid.UUID `db:"instance_id"`
	ZoneId       string    `db:"zone_id"`
	DefinitionId string    `db:"definition_id"`
}

func getPlayerInventoryData(playerName string) (result []PlayerInventoryData, err error) {
	log.Printf("DB - Getting PlayerInventory data for %s", playerName)
	result = []PlayerInventoryData{}
	err = watchdb.Select(&result, "SELECT player_name, instance_id, zone_id, definition_id FROM player_inventory WHERE player_name = $1", playerName)
	return
}

func savePlayerInventory(tx *sqlx.Tx, player player.Player) (err error) {
	log.Printf("DB - Saving Player Inventory for player %s", player.GetName())

	// TODO
	return nil
}
