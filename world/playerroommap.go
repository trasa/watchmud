package world

import (
	"github.com/trasa/watchmud/player"
	"github.com/trasa/watchmud/spaces"
)

type playerRoomMap struct {
	playerToRoom map[*player.Player]*spaces.Room
}

func newPlayerRoomMap() *playerRoomMap {
	return &playerRoomMap{
		playerToRoom: make(map[*player.Player]*spaces.Room),
	}
}

func (m *playerRoomMap) Update(p *player.Player, r *spaces.Room) {
	m.playerToRoom[p] = r
}

func (m *playerRoomMap) Remove(p *player.Player) {
	delete(m.playerToRoom, p)
}

func (m *playerRoomMap) Get(p *player.Player) *spaces.Room {
	return m.playerToRoom[p]
}
