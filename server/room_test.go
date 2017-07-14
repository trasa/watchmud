package server

import "testing"

func TestAddPlayerToRoom(t *testing.T) {
	room := NewRoom(nil, "id", "name", "description")
	player := Player{
		Name: "foo",
	}

	room.AddPlayer(&player)

	if player.Room != room {
		t.Error("player.Room is wrong")
	}
	if room.Players["foo"] == nil {
		t.Error("room.Players is wrong")
	}
}
