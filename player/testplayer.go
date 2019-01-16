package player

import (
	"github.com/trasa/watchmud/combat"
	"github.com/trasa/watchmud/object"
	"log"
)

type TestPlayer struct {
	id        int64
	name      string
	sent      []interface{}
	inventory *PlayerInventory
	slots     *object.Slots
	curHealth int64
	maxHealth int64
	dirty     bool
}

// create a new test player that can track sent messages through 'sentmessages'
func NewTestPlayer(name string) *TestPlayer {
	p := &TestPlayer{
		name:      name,
		inventory: NewPlayerInventory(),
		slots:     object.NewSlots(),
		curHealth: 100,
		maxHealth: 100,
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

func (p *TestPlayer) GetId() int64 {
	return p.id
}

func (p *TestPlayer) GetName() string {
	return p.name
}

func (p *TestPlayer) GetInventory() *PlayerInventory {
	return p.inventory

}

func (p TestPlayer) Slots() *object.Slots {
	return p.slots
}

func (p *TestPlayer) GetCurrentHealth() int64 {
	return p.curHealth
}

func (p *TestPlayer) GetMaxHealth() int64 {
	return p.maxHealth
}

func (p *TestPlayer) TakeMeleeDamage(damage int64) (isDead bool) {
	p.dirty = true
	p.curHealth = p.curHealth - damage
	return p.curHealth <= 0
}

func (p *TestPlayer) IsDead() bool {
	return p.curHealth > 0 // TODO
}

func (p *TestPlayer) CombatantType() combat.CombatantType {
	return combat.PlayerCombatant
}

func (p *TestPlayer) ResetDirtyFlag() {
	p.dirty = false
	p.slots.ResetDirtyFlag()
}

func (p *TestPlayer) IsDirty() bool {
	return p.dirty
}
