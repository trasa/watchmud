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

	for _, a := range player.GetInventory().GetAdded() {
		_, err = tx.NamedExec("INSERT INTO player_inventory (player_name, instance_id, zone_id, definition_id) VALUES (:player_name, :instance_id, :zone_id, :definition_id)",
			map[string]interface{}{
				"player_name":   player.GetName(),
				"instance_id":   a.InstanceId,
				"zone_id":       a.Definition.ZoneId,
				"definition_id": a.Definition.Identifier(),
			})
		if err != nil {
			return
		}
	}

	for _, d := range player.GetInventory().GetRemoved() {
		_, err = tx.NamedExec("DELETE FROM player_inventory WHERE player_name = :player_name AND instance_id = :instance_id",
			map[string]interface{}{
				"player_name": player.GetName(),
				"instance_id": d.InstanceId,
			})
		if err != nil {
			return
		}
	}
	return nil
}
