package player

import (
	"iter"

	"github.com/trasa/watchmud/ordered"
)

// List of players, in the order they joined, keyed by NameKey so a name
// typed in any case finds them.
//
// Add and Remove don't report the duplicate/missing cases the underlying
// ordered.List distinguishes: no caller could do anything about a player who
// is somehow in the room twice, and every one of them is a void method deep
// in a move. They log instead.
type List struct {
	players *ordered.List[string, *Player]
}

func NewList() *List {
	return &List{
		players: ordered.NewList(func(p *Player) string { return NameKey(p.Name()) }),
	}
}

func (l *List) Add(p *Player) {
	if err := l.players.Add(p); err != nil {
		p.Log().Warn().Err(err).Msg("player.List.Add")
	}
}

func (l *List) Remove(p *Player) {
	if err := l.players.Remove(p); err != nil {
		p.Log().Warn().Err(err).Msg("player.List.Remove")
	}
}

// All the players, in the order they joined.
func (l *List) All() iter.Seq[*Player] {
	return l.players.All()
}

// AllExcept every player but this one, in the order they joined. A nil
// exclusion excludes nobody, which is how an unattributed room description
// asks for the whole room.
func (l *List) AllExcept(exclude *Player) iter.Seq[*Player] {
	if exclude == nil {
		return l.All()
	}
	return l.players.AllExcept(NameKey(exclude.Name()))
}

// Slice of the players, in the order they joined. Prefer All; this is for
// callers that index or hold on to the result.
func (l *List) Slice() []*Player {
	return l.players.Slice()
}

func (l *List) FindByName(name string) *Player {
	p, _ := l.players.Get(NameKey(name)) // nil when nobody has it
	return p
}

func (l *List) Count() int {
	return l.players.Len()
}
