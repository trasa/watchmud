package world

import (
	"fmt"
	"log"
)

var knownPlayersById = make(map[string]*Player)

func ClearKnownPlayers() {
	knownPlayersById = make(map[string]*Player)
}

func GetAllPlayers() (result []*Player) {
	for _, v := range knownPlayersById {
		result = append(result, v)
	}
	return result
}

func AddKnownPlayer(player *Player) {
	knownPlayersById[player.Id] = player
}

func Login(player *Player) error {
	if player.Id == "" {
		log.Printf("Player does not have a valid id: %s", player)
		return fmt.Errorf("Login: ID is not valid")
	}
	log.Printf("%s logging in", player)
	AddKnownPlayer(player)

	// put the player in the start room
	worldInstance.StartRoom.AddPlayer(player)
	return nil
}

func Logout(playerId string) error {
	log.Printf("%s logged out", playerId)
	player := knownPlayersById[playerId]
	player.Room.RemovePlayer(player)
	delete(knownPlayersById, playerId)
	return nil
}
