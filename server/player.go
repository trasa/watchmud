package server

import (
	"fmt"
	"log"
)

// TODO move to world
var playersByClient = make(map[*Client]*Player)

type Player struct {
	Name   string
	Room   *Room   `json:"-"`
	Client *Client `json:"-"`
}

func NewPlayer(name string, client *Client) *Player {
	p := Player{
		Name:   name,
		Client: client,
	}
	return &p
}

func (p *Player) String() string {
	return fmt.Sprintf("(Player Name='%s' in room '%v')", p.Name, p.Room)
}

func (p *Player) FindZone() *Zone {
	if p.Room != nil {
		return p.Room.Zone
	}
	// TODO return err?
	return nil
}

func (p *Player) Send(message interface{}) {
	if p.Client != nil {
		p.Client.source <- message
	} else {
		log.Printf("Can't send message to player: no client attached: %v", message)
		// TODO return err
	}
}

// Locate the player who is attached to this client
// TODO move to world
func FindPlayerByClient(c *Client) *Player {
	return playersByClient[c]
}
