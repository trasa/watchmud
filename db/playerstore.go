package db

type PlayerData struct {
	Name string `db:"player_name"`
	//inventory *player.PlayerInventory // join to inventory table, instance ids?
	//slots     *object.Slots // maps instance ids to location
	CurHealth int64 `db:"current_health"`
	MaxHealth int64 `db:"max_health"`
	// current location of the player? "zoneId.roomId" ?
}

func GetPlayerData(playerName string) (result *PlayerData, err error) {
	result = &PlayerData{}
	err = watchdb.Get(result, "SELECT player_name, current_health, max_health FROM players where player_name = $1", playerName)
	return
}
