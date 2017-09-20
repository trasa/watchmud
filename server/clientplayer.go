package server

import (
	"fmt"
	"github.com/trasa/watchmud/client"
	"github.com/trasa/watchmud/player"
	"github.com/trasa/watchmud/testclient"
)

type ClientPlayer struct {
	Name      string
	Client    client.Client `json:"-"`
	Inventory map[string][]player.InventoryItem
}

// Create a ClientPlayer connected to a new TestClient
// (for testing)
func NewTestClientPlayer(name string) (p *ClientPlayer, cli *testclient.TestClient) {
	p = NewClientPlayer(name, nil)
	cli = testclient.NewTestClient(p)
	p.Client = cli
	return
}

// Create a new player and set it up to work with this client
func NewClientPlayer(name string, client client.Client) *ClientPlayer {
	p := ClientPlayer{
		Name:      name,
		Client:    client, // address of interface
		Inventory: make(map[string][]player.InventoryItem),
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

func (p *ClientPlayer) GetInventory() map[string][]player.InventoryItem {
	return p.Inventory
}

func (p *ClientPlayer) AddInventory(item player.InventoryItem) {
	if val, ok := p.Inventory[item.Id]; ok {
		val = append(val, item)
	} else {
		p.Inventory[item.Id] = []player.InventoryItem{item}
	}
}
