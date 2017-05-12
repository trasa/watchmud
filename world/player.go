package world

import (
	"fmt"
)

type Player struct {
	Id   string
	Name string
	Room *Room `json:"-"`
}

func NewPlayer(id string, name string) *Player {
	p := Player{
		Name: name,
		Id:   id,
	}
	return &p
}

func (p *Player) String() string {
	return fmt.Sprintf("(Player Id='%s', Name='%s' in room '%v')", p.Id, p.Name, p.Room)
}

func (p *Player) FindZone() *Zone {
	if p.Room != nil {
		return p.Room.Zone
	}
	return nil
}
