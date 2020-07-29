package db

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"github.com/satori/go.uuid"
	"github.com/trasa/watchmud/player"
)

type PlayerInventoryData struct {
	InventoryId  int32     `db:"inventory_id"`
	PlayerId     int64     `db:"player_id"`
	InstanceId   uuid.UUID `db:"instance_id"`
	ZoneId       string    `db:"zone_id"`
	DefinitionId string    `db:"definition_id"`
}

func getPlayerInventoryData(playerId int64) (result []PlayerInventoryData, err error) {
	result = []PlayerInventoryData{}
	err = watchdb.Select(&result, "SELECT inventory_id, player_id, instance_id, zone_id, definition_id FROM player_inventory WHERE player_id = $1 ORDER BY inventory_id", playerId)
	return
}

func savePlayerInventory(tx *sqlx.Tx, player player.Player) (err error) {
	log.Printf("DB - Saving Player Inventory for player %s %d", player.GetName(), player.GetId())

	for _, a := range player.Inventory().GetAdded() {
		_, err = tx.NamedExec("INSERT INTO player_inventory (player_id, instance_id, zone_id, definition_id) VALUES (:player_id, :instance_id, :zone_id, :definition_id)",
			map[string]interface{}{
				"player_id":     player.GetId(),
				"instance_id":   a.InstanceId,
				"zone_id":       a.Definition.ZoneId,
				"definition_id": a.Definition.Identifier(),
			})
		if err != nil {
			return
		}
	}

	for _, d := range player.Inventory().GetRemoved() {
		_, err = tx.NamedExec("DELETE FROM player_inventory WHERE player_id = :player_id AND instance_id = :instance_id",
			map[string]interface{}{
				"player_id":   player.GetId(),
				"instance_id": d.InstanceId,
			})
		if err != nil {
			return
		}
	}
	return nil
}
