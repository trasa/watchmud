package server

import (
	"fmt"
	"log"
)

// TODO move to world
var playersByClient = make(map[*Client]*Player)

// Locate the player who is attached to this client
// TODO move to world
func FindPlayerByClient(c *Client) *Player {
	return playersByClient[c]
}

type Player struct {
	Name   string
	Room   *Room        `json:"-"`
	Client *Client      `json:"-"`
	Send   playerSender `json:"-"`
}

// method called to send a message to the player
// could be overridden for test purposes
type playerSender func(message interface{})

// Create a new player and set it up to work with this client
func NewPlayer(name string, client *Client) *Player {
	p := Player{
		Name:   name,
		Client: client,
	}
	p.Send = func(message interface{}) {
		log.Printf("using REAL playersender to talk to %s", p.Name)
		if p.Client != nil {
			p.Client.source <- message
		} else {
			log.Printf("Can't send message to player: no client attached: %v", message)
			// TODO return err
		}
	}
	return &p
}

func (p *Player) String() string {
	return fmt.Sprintf("(Player Name='%s' in room '%v')", p.Name, p.Room)
}

// TODO move to somewhere else?
func (p *Player) FindZone() *Zone {
	if p.Room != nil {
		return p.Room.Zone
	}
	// TODO return err?
	return nil
}
