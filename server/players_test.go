package server

import "testing"

func TestPlayers_Add(t *testing.T) {
	players := newPlayers()
	p := &Player{}

	players.Add(p)
	if _, ok := players.players[p]; !ok {
		t.Error("Player added to Players but not found by key")
	}
}

func TestPlayers_Remove(t *testing.T) {
	players := newPlayers()
	p := &Player{}

	players.Add(p)
	players.Remove(p)

	if _, ok := players.players[p]; ok {
		t.Error("Player removed from players but it is still there")
	}
}

func TestPlayers_Iter(t *testing.T) {
	players := newPlayers()
	p := &Player{}

	players.Add(p)
	count := 0
	players.Iter(func(p *Player) {
		count++
	})
	if count != 1 {
		t.Errorf("found %d expected 1", count)
	}

}
