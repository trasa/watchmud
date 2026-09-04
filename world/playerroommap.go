package world

import (
	"github.com/trasa/syncmap"
	"github.com/trasa/watchmud/player"
	"github.com/trasa/watchmud/spaces"
)

type PlayerRoomMap struct {
	playerToRoom  map[*player.Player]*spaces.Room
	roomToPlayers syncmap.MapList // TODO replace this
}

func NewPlayerRoomMap() *PlayerRoomMap {
	return &PlayerRoomMap{
		playerToRoom:  make(map[*player.Player]*spaces.Room),
		roomToPlayers: syncmap.NewMapList(),
	}
}

func (m *PlayerRoomMap) Add(p *player.Player, r *spaces.Room) {
	m.playerToRoom[p] = r
	m.roomToPlayers.Add(r, p)
}

func (m *PlayerRoomMap) Remove(p *player.Player) {
	r := m.playerToRoom[p]
	delete(m.playerToRoom, p)
	if r != nil {
		m.roomToPlayers.RemoveItem(r, p)
		r.RemovePlayer(p)
	}
}

func (m *PlayerRoomMap) Get(p *player.Player) *spaces.Room {
	return m.playerToRoom[p]
}

func (m *PlayerRoomMap) GetPlayers(r *spaces.Room) []*player.Player {
	players := make([]*player.Player, 0, len(m.roomToPlayers.Get(r)))
	for _, p := range m.roomToPlayers.Get(r) {
		players = append(players, p.(*player.Player))
	}
	return players
}
