package world

import (
	"github.com/trasa/syncmap"
	"github.com/trasa/watchmud/player"
	"github.com/trasa/watchmud/spaces"
	"sync"
)

type PlayerRoomMap struct {
	sync.RWMutex
	playerToRoom  map[player.Player]*spaces.Room
	roomToPlayers syncmap.MapList
}

func NewPlayerRoomMap() *PlayerRoomMap {
	return &PlayerRoomMap{
		playerToRoom:  make(map[player.Player]*spaces.Room),
		roomToPlayers: syncmap.NewMapList(),
	}
}

func (m *PlayerRoomMap) Add(p player.Player, r *spaces.Room) {
	m.Lock()
	defer m.Unlock()
	m.playerToRoom[p] = r
	m.roomToPlayers.Add(r, p)
}

func (m *PlayerRoomMap) Remove(p player.Player) {
	m.Lock()
	defer m.Unlock()
	r := m.playerToRoom[p]
	delete(m.playerToRoom, p)
	if r != nil {
		m.roomToPlayers.RemoveItem(r, p)
		r.PlayerList.Remove(p)
	}
}

func (m *PlayerRoomMap) Get(p player.Player) *spaces.Room {
	m.RLock()
	defer m.RUnlock()
	return m.playerToRoom[p]
}

func (m *PlayerRoomMap) GetPlayers(r *spaces.Room) []player.Player {
	m.RLock()
	defer m.RUnlock()

	players := make([]player.Player, 0, len(m.roomToPlayers.Get(r)))
	for _, p := range m.roomToPlayers.Get(r) {
		players = append(players, p.(player.Player))
	}
	return players
}
