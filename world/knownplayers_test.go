package world

import (
	"log"
	"testing"
)

func TestLogin(t *testing.T) {
	ClearKnownPlayers()
	originalPlayer := NewPlayer("id", "Name")
	if err := Login(originalPlayer); err != nil {
		t.Error("Failed to login")
	}
}

func TestGetAll(t *testing.T) {
	ClearKnownPlayers()
	playerA := NewPlayer("a", "a")
	playerB := NewPlayer("b", "b")
	if err := Login(playerA); err != nil {
		t.Error("Failed to login A")
	}
	if err := Login(playerB); err != nil {
		t.Error("Failed to login B")
	}

	players := GetAllPlayers()
	if len(players) != 2 {
		t.Errorf("Expected 2 players, found %d", len(players))
	}
	for _, p := range players {
		log.Printf("Player: %s", p)
	}
}
