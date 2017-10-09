package server

import (
	"fmt"
	"github.com/trasa/watchmud/client"
	"github.com/trasa/watchmud/object"
	"github.com/trasa/watchmud/thing"
)

type ClientPlayer struct {
	Name      string
	Client    client.Client `json:"-"`
	Inventory thing.Map
}

// Create a ClientPlayer connected to a new TestClient
// (for testing)
func NewTestClientPlayer(name string) (p *ClientPlayer, cli *client.TestClient) {
	p = NewClientPlayer(name, nil)
	cli = client.NewTestClient(p)
	p.Client = cli
	return
}

// Create a new player and set it up to work with this client
func NewClientPlayer(name string, client client.Client) *ClientPlayer {
	p := ClientPlayer{
		Name:      name,
		Client:    client, // address of interface
		Inventory: make(thing.Map),
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
	return fmt.Sprintf("(Player Name='%s')", p.Name)
}

func (p *ClientPlayer) GetInventoryMap() thing.Map {
	return p.Inventory
}

func (p *ClientPlayer) AddInventory(instance *object.Instance) error {
	return p.Inventory.Add(instance)
}

func (p *ClientPlayer) RemoveInventory(instance *object.Instance) error {
	return p.Inventory.Remove(instance)
}
