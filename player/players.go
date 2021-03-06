package player

import "sync"

type List struct {
	sync.RWMutex
	players map[Player]Player
	byName  map[string]Player
}

func NewList() *List {
	return &List{
		players: make(map[Player]Player),
		byName:  make(map[string]Player),
	}
}

func (ps *List) Add(p Player) {
	ps.Lock()
	defer ps.Unlock()
	ps.players[p] = p
	ps.byName[p.GetName()] = p
}

func (ps *List) Remove(p Player) {
	ps.Lock()
	defer ps.Unlock()
	delete(ps.players, p)
	delete(ps.byName, p.GetName())
}

func (ps *List) GetAll() []Player {
	ps.RLock()
	defer ps.RUnlock()
	// copy the keys into a new slice
	// and return that slice
	keys := []Player{}
	for p := range ps.players {
		keys = append(keys, p)
	}
	return keys
}

func (ps *List) Iter(routine func(Player)) {
	ps.RLock()
	defer ps.RUnlock()
	for p := range ps.players {
		routine(p)
	}
}

func (ps *List) FindByName(name string) Player {
	ps.RLock()
	defer ps.RUnlock()
	return ps.byName[name]
}

func (ps *List) Count() int {
	ps.RLock()
	defer ps.RUnlock()
	return len(ps.players)
}
