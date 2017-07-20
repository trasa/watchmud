package server

import (
	"fmt"
	"github.com/trasa/watchmud/client"
)

type ClientPlayer struct {
	Name   string
	Client client.Client `json:"-"`
}

// Create a new player and set it up to work with this client
func NewClientPlayer(name string, client client.Client) *ClientPlayer {
	p := ClientPlayer{
		Name:   name,
		Client: client, // address of interface
	}
	return &p
}

func (p *ClientPlayer) GetName() string {
	return p.Name
}

func (p *ClientPlayer) Send(msg interface{}) {
	p.Client.Send(msg)
}

func (p *ClientPlayer) String() string {
	return fmt.Sprintf("(Player Name='%s')", p.Name )
}
//
//// TODO move to somewhere else?
//func (p *ClientPlayer) FindZone() *Zone {
//	if p.Room != nil {
//		return p.Room.Zone
//	}
//	// TODO return err?
//	return nil
//}
