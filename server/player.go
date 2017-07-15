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
	Name         string
	Room         *Room   `json:"-"`
	Client       *Client `json:"-"`
	PlayerSender `json:"-"`
}

type PlayerSender func(player *Player, message interface{})

func NewPlayer(name string, client *Client) *Player {
	p := Player{
		Name:         name,
		Client:       client,
		PlayerSender: sendToPlayer, // use real method that relies on client
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
	p.PlayerSender(p, message)
}

// Really send a message to the player via the client
func sendToPlayer(p *Player, message interface{}) {
	if p.Client != nil {
		p.Client.source <- message
	} else {
		log.Printf("Can't send message to player: no client attached: %v", message)
		// TODO return err
	}
}
