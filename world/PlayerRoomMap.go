package world

import (
	"github.com/trasa/watchmud/player"
	"sync"
)

type PlayerRoomMap struct {
	sync.RWMutex
	playerToRoom map[player.Player]*Room
}

func NewPlayerRoomMap() *PlayerRoomMap {
	return &PlayerRoomMap{
		playerToRoom: make(map[player.Player]*Room),
	}
}

func (m *PlayerRoomMap) Add(p player.Player, r *Room) {
	m.Lock()
	defer m.Unlock()
	m.playerToRoom[p] = r
}

func (m *PlayerRoomMap) Remove(p player.Player) {
	m.Lock()
	defer m.Unlock()
	delete(m.playerToRoom, p)
}

func (m *PlayerRoomMap) Get(p player.Player) *Room {
	m.RLock()
	defer m.RUnlock()
	return m.playerToRoom[p]
}
