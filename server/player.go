package server

import (
	"fmt"
	"log"
)

var playersByClient = make(map[*Client]*Player)

type Player struct {
	Id     string
	Name   string
	Room   *Room   `json:"-"`
	Client *Client `json:"-"`
}

func NewPlayer(id string, name string, client *Client) *Player {
	p := Player{
		Name:   name,
		Id:     id,
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

func (p *Player) Send(message interface{}) {
	if p.Client != nil {
		p.Client.source <- message
	} else {
		log.Printf("Can't send message to player: no client attached: %v", message)
	}
}

// Locate the player who is attached to this client
func FindPlayerByClient(c *Client) *Player {
	return playersByClient[c]
}
