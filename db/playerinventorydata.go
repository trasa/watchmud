package db

import "github.com/satori/go.uuid"

type PlayerInventoryData struct {
	PlayerName   string    `db:"player_name"`
	InstanceId   uuid.UUID `db:"instance_id"`
	ZoneId       string    `db:"zone_id"`
	DefinitionId string    `db:"definition_id"`
}

func GetPlayerInventoryData(playerName string) (result []PlayerInventoryData, err error) {
	result = []PlayerInventoryData{}
	err = watchdb.Select(&result, "SELECT player_name, instance_id, zone_id, definition_id FROM player_inventory WHERE player_name = $1", playerName)
	return
}
