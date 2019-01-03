package server

import (
	"fmt"
	"github.com/trasa/watchmud/client"
	"github.com/trasa/watchmud/combat"
	"github.com/trasa/watchmud/db"
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
	race      int32
	class     int32
	dirty     bool
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
func NewClientPlayer(name string, client client.Client) *ClientPlayer { // TODO need race and class
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

// Load player information into this struct without flagging anything as dirty
func (p *ClientPlayer) LoadPlayerData(pd *db.PlayerData) {
	p.Name = pd.Name
	p.curHealth = pd.CurHealth
	p.maxHealth = pd.MaxHealth
	p.race = pd.Race
	p.class = pd.Class

	// TODO: other stuff
}

// Load player inventory into this struct without flagging anything as dirty
func (p *ClientPlayer) LoadInventory(instance *object.Instance) {
	p.inventory.Load(instance)
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

func (p *ClientPlayer) GetInventory() *player.PlayerInventory {
	return p.inventory
}

//func (p *ClientPlayer) GetInventoryById(id uuid.UUID) (*object.Instance, bool) {
//	return p.inventory.GetByInstanceId(id)
//}

//func (p *ClientPlayer) GetInventoryByName(name string) (*object.Instance, bool) {
//	return p.inventory.GetByName(name)
//}

//func (p *ClientPlayer) GetAllInventory() []*object.Instance {
//	return p.inventory.GetAll()
//}

//func (p *ClientPlayer) AddInventory(instance *object.Instance) error {
//	p.dirty = true
//	return p.inventory.Add(instance)
//}

//func (p *ClientPlayer) RemoveInventory(instance *object.Instance) error {
//	p.dirty = true
//	return p.inventory.Remove(instance)
//}

//func (p *ClientPlayer) FindInventory(target string) (*object.Instance, bool) {
//	return p.inventory.Find(target)
//}

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
	p.dirty = true
	p.curHealth = p.curHealth - damage
	return p.curHealth <= 0
}

func (p *ClientPlayer) IsDead() bool {
	return p.curHealth <= 0
}

func (p *ClientPlayer) CombatantType() combat.CombatantType {
	return combat.PlayerCombatant
}

func (p *ClientPlayer) ResetDirtyFlag() {
	p.dirty = false
	p.slots.ResetDirtyFlag()
}

func (p *ClientPlayer) IsDirty() bool {
	return p.dirty || p.slots.IsDirty()
}
