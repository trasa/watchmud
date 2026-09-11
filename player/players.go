package player

// List of players
// TODO replace with generic
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

func (l *List) Add(p *Player) {
	l.players[p] = p
	l.byName[p.name] = p
}

func (l *List) Remove(p *Player) {
	delete(l.players, p)
	delete(l.byName, p.name)
}

func (l *List) GetAll() []*Player {
	// copy the keys into a new slice
	// and return that slice
	var keys []*Player
	for p := range l.players {
		keys = append(keys, p)
	}
	return keys
}

func (l *List) GetExcept(exclude *Player) []*Player {
	var result []*Player
	for p := range l.players {
		if exclude != p {
			result = append(result, p)
		}
	}
	return result
}

func (l *List) Iter(routine func(*Player)) {
	for p := range l.players {
		routine(p)
	}
}

func (l *List) FindByName(name string) *Player {
	// TODO what happens if name is not found?
	return l.byName[name]
}

func (l *List) Count() int {
	// TODO replace with support for len
	return len(l.players)
}
