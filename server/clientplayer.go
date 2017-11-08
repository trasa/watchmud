package server

import (
	"fmt"
	"github.com/satori/go.uuid"
	"github.com/trasa/watchmud/client"
	"github.com/trasa/watchmud/object"
	"github.com/trasa/watchmud/thing"
)

type ClientPlayer struct {
	Name      string
	Client    client.Client `json:"-"`
	inventory thing.Map     // TODO replace with something better...
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
		inventory: make(thing.Map),
	}
	return &p
}

func (p *ClientPlayer) GetName() string {
	return p.Name
}

func (p *ClientPlayer) Send(msg interface{}) error {
	return p.Client.Send(msg)
}

func (p *ClientPlayer) String() string {
	return fmt.Sprintf("(Player Name='%s')", p.Name)
}

func (p *ClientPlayer) GetInventoryById(id uuid.UUID) (inst *object.Instance, exists bool) {
	t, exists := p.inventory.Get(id.String())
	if exists {
		inst = t.(*object.Instance)
	}
	return
}

func (p *ClientPlayer) GetInventoryByName(name string) (inst *object.Instance, exists bool) {
	for _, t := range p.inventory {
		if name == t.(*object.Instance).Definition.Name {
			return t.(*object.Instance), true
		}
	}
	return nil, false
}

func (p *ClientPlayer) GetAllInventory() (result []*object.Instance) {
	for _, t := range p.inventory {
		result = append(result, t.(*object.Instance))
	}
	return result
}

func (p *ClientPlayer) AddInventory(instance *object.Instance) error {
	return p.inventory.Add(instance)
}

func (p *ClientPlayer) RemoveInventory(instance *object.Instance) error {
	return p.inventory.Remove(instance)
}
