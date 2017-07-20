package server

import "sync"

type Players struct {
	sync.Mutex
	players map[*Player]*Player
}

func newPlayers() *Players {
	return &Players{
		players: make(map[*Player]*Player),
	}
}

func (ps *Players) Add(p *Player) {
	ps.Lock()
	defer ps.Unlock()
	ps.players[p] = p
}

func (ps *Players) Remove(p *Player) {
	ps.Lock()
	defer ps.Unlock()
	delete(ps.players, p)
}

func (ps *Players) Iter(routine func(*Player)) {
	ps.Lock()
	defer ps.Unlock()
	for p := range ps.players {
		routine(p)
	}
}
