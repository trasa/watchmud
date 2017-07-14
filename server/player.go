package server

import (
	"fmt"
)
var playersByClient = make(map[*Client]*Player)

type Player struct {
	Id   string
	Name string
	Room *Room `json:"-"`
	Client *Client `json:"-"`
}

func NewPlayer(id string, name string, client *Client) *Player {
	p := Player{
		Name: name,
		Id:   id,
		Client: client,
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

// Locate the player who is attached to this client
func FindPlayerByClient(c *Client) *Player {
	return playersByClient[c];
}

