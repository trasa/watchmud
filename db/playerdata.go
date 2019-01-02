package db

import "log"

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
	result = &PlayerData{}
	err = watchdb.Get(result, "SELECT player_name, current_health, max_health, race, class FROM players where player_name = $1", playerName)
	if err != nil {
		return
	}
	result.Inventory, err = GetPlayerInventoryData(playerName)
	if err != nil {
		log.Printf("err %v", err)
	}

	for _, inv := range result.Inventory {
		log.Printf("inv %v", inv)
	}
	return
}

func CreatePlayerData(playerName string, race int32, class int32) (result *PlayerData, err error) {
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
	return
}
