package player

import (
	"github.com/satori/go.uuid"
	"github.com/trasa/watchmud/object"
	"github.com/trasa/watchmud/thing"
	"log"
)

type TestPlayer struct {
	name      string
	sent      []interface{}
	inventory thing.Map
	slots     *object.Slots
}

// create a new test player that can track sent messages through 'sentmessages'
func NewTestPlayer(name string) *TestPlayer {
	p := &TestPlayer{
		name:      name,
		inventory: make(thing.Map),
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

func (p *TestPlayer) GetInventoryById(id uuid.UUID) (inst *object.Instance, exists bool) {
	t, exists := p.inventory[id.String()]
	if exists {
		inst = t.(*object.Instance)
	}
	return
}

func (p *TestPlayer) GetInventoryByName(name string) (*object.Instance, bool) {
	for _, t := range p.inventory {
		if name == t.(*object.Instance).Definition.Name {
			return t.(*object.Instance), true
		}
	}
	return nil, false
}

func (p *TestPlayer) GetAllInventory() (result []*object.Instance) {
	for _, t := range p.inventory {
		result = append(result, t.(*object.Instance))
	}
	return
}

func (p *TestPlayer) Inventory() thing.Map {
	return p.inventory
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
