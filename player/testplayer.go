package player

import (
	"github.com/trasa/watchmud/object"
	"github.com/trasa/watchmud/thing"
	"log"
)

type TestPlayer struct {
	name      string
	sent      []interface{}
	inventory thing.Map
}

// create a new test player that can track sent messages through 'sentmessages'
func NewTestPlayer(name string) *TestPlayer {
	p := &TestPlayer{
		name:      name,
		inventory: make(thing.Map),
	}
	return p
}

func (p *TestPlayer) Send(msg interface{}) {
	log.Printf("sending for player %s", msg)
	p.sent = append(p.sent, msg)
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

func (p *TestPlayer) GetInventoryMap() thing.Map {
	return p.inventory
}

func (p *TestPlayer) AddInventory(instance *object.Instance) error {
	return p.inventory.Add(instance)
}

func (p *TestPlayer) RemoveInventory(instance *object.Instance) error {
	return p.inventory.Remove(instance)
}
