package server

import (
	"fmt"
	"github.com/satori/go.uuid"
	"github.com/trasa/watchmud/client"
	"github.com/trasa/watchmud/combat"
	"github.com/trasa/watchmud/object"
	"github.com/trasa/watchmud/player"
)

type ClientPlayer struct {
	Name      string
	Client    client.Client
	inventory *player.PlayerInventory
	slots     *object.Slots
	curHealth int64
	maxHealth int64
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
		inventory: player.NewPlayerInventory(),
		slots:     object.NewSlots(),
		curHealth: 100, // TODO need a default health here
		maxHealth: 100,
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

func (p *ClientPlayer) GetInventoryById(id uuid.UUID) (*object.Instance, bool) {
	return p.inventory.GetByInstanceId(id)
}

func (p *ClientPlayer) GetInventoryByName(name string) (*object.Instance, bool) {
	return p.inventory.GetByName(name)
}

func (p *ClientPlayer) GetAllInventory() []*object.Instance {
	return p.inventory.GetAll()
}

func (p *ClientPlayer) AddInventory(instance *object.Instance) error {
	return p.inventory.Add(instance)
}

func (p *ClientPlayer) RemoveInventory(instance *object.Instance) error {
	return p.inventory.Remove(instance)
}

func (p *ClientPlayer) FindInventory(target string) (*object.Instance, bool) {
	return p.inventory.Find(target)
}

func (p *ClientPlayer) Slots() *object.Slots {
	return p.slots
}

func (p *ClientPlayer) GetCurrentHealth() int64 {
	return p.curHealth
}

func (p *ClientPlayer) GetMaxHealth() int64 {
	return p.maxHealth
}

func (p *ClientPlayer) TakeMeleeDamage(damage int64) (isDead bool) {
	p.curHealth = p.curHealth - damage
	return p.curHealth <= 0
}

func (p *ClientPlayer) IsDead() bool {
	return p.curHealth <= 0
}

func (p *ClientPlayer) CombatantType() combat.CombatantType {
	return combat.PlayerCombatant
}
