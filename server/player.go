package server

import (
	"fmt"
)

// TODO move to world
var playersByClient = make(map[Client]*Player)

// Locate the player who is attached to this client
// TODO move to world
func FindPlayerByClient(c Client) *Player {
	return playersByClient[c]
}

type Player struct {
	Name   string
	Room   *Room  `json:"-"`
	Client Client `json:"-"`
}

// Create a new player and set it up to work with this client
func NewPlayer(name string, client Client) *Player {
	p := Player{
		Name:   name,
		Client: client, // address of interface
	}
	return &p
}

func (p *Player) Send(msg interface{}) {
	p.Client.Send(msg)
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
