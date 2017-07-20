package player

import "sync"

type Players struct {
	sync.Mutex
	players map[Player]Player
}

func NewPlayers() *Players {
	return &Players{
		players: make(map[Player]Player),
	}
}

func (ps *Players) Add(p Player) {
	ps.Lock()
	defer ps.Unlock()
	ps.players[p] = p
}

func (ps *Players) Remove(p Player) {
	ps.Lock()
	defer ps.Unlock()
	delete(ps.players, p)
}

func (ps *Players) All() []Player {
	ps.Lock()
	defer ps.Unlock()
	// copy the keys into a new slice
	// and return that slice
	keys := []Player{}
	for p := range ps.players {
		keys = append(keys, p)
	}
	return keys
}

func (ps *Players) Iter(routine func(Player)) {
	ps.Lock()
	defer ps.Unlock()
	for p := range ps.players {
		routine(p)
	}
}
