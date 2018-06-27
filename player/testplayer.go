package player

import (
	"github.com/satori/go.uuid"
	"github.com/trasa/watchmud/object"
	"log"
)

type TestPlayer struct {
	name      string
	sent      []interface{}
	inventory *PlayerInventory
	slots     *object.Slots
	curHealth int
	maxHealth int
}

// create a new test player that can track sent messages through 'sentmessages'
func NewTestPlayer(name string) *TestPlayer {
	p := &TestPlayer{
		name:      name,
		inventory: NewPlayerInventory(),
		slots:     object.NewSlots(),
	}
	return p
}

func (p *TestPlayer) Send(msg interface{}) error {
	log.Printf("sending for player %s", msg)
	p.sent = append(p.sent, msg)
	return nil
}

func (p *TestPlayer) SentMessageCount() int {
	return len(p.sent)
}

func (p *TestPlayer) GetSentResponse(i int) interface{} {
	return p.sent[i]
}

func (p *TestPlayer) GetName() string {
	return p.name
}

func (p *TestPlayer) GetInventoryById(id uuid.UUID) (*object.Instance, bool) {
	return p.inventory.GetByInstanceId(id)
}

func (p *TestPlayer) GetInventoryByName(name string) (*object.Instance, bool) {
	return p.inventory.GetByName(name)
}

func (p *TestPlayer) FindInventory(target string) (*object.Instance, bool) {
	return p.inventory.Find(target)
}

func (p *TestPlayer) GetAllInventory() []*object.Instance {
	return p.inventory.GetAll()
}

func (p *TestPlayer) AddInventory(instance *object.Instance) error {
	return p.inventory.Add(instance)
}

func (p *TestPlayer) RemoveInventory(instance *object.Instance) error {
	return p.inventory.Remove(instance)
}

func (p TestPlayer) Slots() *object.Slots {
	return p.slots
}

func (p *TestPlayer) GetCurrentHealth() int {
	return p.curHealth
}

func (p *TestPlayer) GetMaxHealth() int {
	return p.maxHealth
}
