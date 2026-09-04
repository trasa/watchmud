package player

type List struct {
	players map[*Player]*Player
	byName  map[string]*Player
}

func NewList() *List {
	return &List{
		players: make(map[*Player]*Player),
		byName:  make(map[string]*Player),
	}
}

func (ps *List) Add(p *Player) {
	ps.players[p] = p
	ps.byName[p.Name] = p
}

func (ps *List) Remove(p *Player) {
	delete(ps.players, p)
	delete(ps.byName, p.Name)
}

func (ps *List) GetAll() []*Player {
	// copy the keys into a new slice
	// and return that slice
	var keys []*Player
	for p := range ps.players {
		keys = append(keys, p)
	}
	return keys
}

func (ps *List) Iter(routine func(*Player)) {
	for p := range ps.players {
		routine(p)
	}
}

func (ps *List) FindByName(name string) *Player {
	// TODO what happens if name is not found?
	return ps.byName[name]
}

func (ps *List) Count() int {
	// TODO replace with support for len
	return len(ps.players)
}
