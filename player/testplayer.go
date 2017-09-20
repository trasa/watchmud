package player

import (
	"log"
)

type TestPlayer struct {
	name      string
	sent      []interface{}
	inventory map[string][]InventoryItem
}

// create a new test player that can track sent messages through 'sentmessages'
func NewTestPlayer(name string) *TestPlayer {
	p := &TestPlayer{
		name:      name,
		inventory: make(map[string][]InventoryItem),
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

func (p *TestPlayer) GetInventory() map[string][]InventoryItem {
	return p.inventory
}

func (p *TestPlayer) AddInventory(item InventoryItem) {
	if val, ok := p.inventory[item.Id]; ok {
		val = append(val, item)
	} else {
		p.inventory[item.Id] = []InventoryItem{item}
	}
}
