package player

import (
	"log"
	"testing"
)

func TestPlayers_Add(t *testing.T) {
	players := NewList()
	p := &TestPlayer{}

	players.Add(p)
	if _, ok := players.players[p]; !ok {
		t.Error("Player added to Players but not found by key")
	}
}

func TestPlayers_Remove(t *testing.T) {
	players := NewList()
	p := &TestPlayer{}

	players.Add(p)
	players.Remove(p)

	if _, ok := players.players[p]; ok {
		t.Error("Player removed from players but it is still there")
	}
}

func TestPlayers_RemoveDoesntExist(t *testing.T) {
	players := NewList()
	p := &TestPlayer{}
	players.Remove(p)
	if _, ok := players.players[p]; ok {
		t.Error("Player never existed in list but now it exists?!")
	}
}

func TestPlayers_Iter(t *testing.T) {
	players := NewList()
	p := &TestPlayer{}

	players.Add(p)
	count := 0
	players.Iter(func(p Player) {
		count++
	})
	if count != 1 {
		t.Errorf("found %d expected 1", count)
	}
}

func TestPlayer_GetAll(t *testing.T) {
	players := NewList()
	p := &TestPlayer{name: "a"}
	players.Add(p)

	all := players.GetAll()
	players.Add(&TestPlayer{name: "b"})
	log.Printf("addr of all: %p", &all)
	if len(all) != 1 {
		t.Error("expected len = 1")
	}

	moreAll := players.GetAll()
	log.Printf("addr of moreAll: %p", &moreAll)
	if len(all) != 1 {
		t.Error("still expected len 1")
	}
	if len(moreAll) != 2 {
		t.Error("expected len = 2")
	}
}
