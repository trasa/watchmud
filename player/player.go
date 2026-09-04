package player

import (
	"github.com/trasa/watchmud/object"
	"github.com/trasa/watchmud/rules"
)

// Sender is anything that can deliver a message to this player's connection
type Sender interface {
	Send(msg interface{}) error
}

type Player struct {
	Id        string
	Name      string
	out       Sender // was: client.Client, via ClientPlayer
	Lineage   *rules.Lineage
	Class     *rules.Class
	inventory *Inventory
	slots     *object.Slots
	curHealth int64
	maxHealth int64
	isDirty   bool // If true, need to write back to the database.
	location  *Location
	abilities rules.Abilities
}

func New(id string,
	name string,
	out Sender,
	lineage *rules.Lineage,
	class *rules.Class,
	abilities rules.Abilities,
) *Player {
	return &Player{
		Id:        id,
		Name:      name,
		out:       out,
		Lineage:   lineage,
		Class:     class,
		inventory: NewInventory(),
		slots:     object.NewSlots(),
		curHealth: 100, // TODO need a default here,
		maxHealth: 100,
		abilities: abilities,
	}
}

// Inventory returns the inventory
func (p *Player) Inventory() *Inventory {
	// TODO is this needed? Should p.Inventory become visible?
	// is needing this call indicating a problem?
	return p.inventory
}

// Slots returns the inventory
func (p *Player) Slots() *object.Slots {
	// TODO is this needed? Should p.Inventory become visible?
	// is needing this call indicating a problem?
	return p.slots
}

// NewTestPlayer that tracks messages
func NewTestPlayer(id string, name string, out Sender) *Player {
	if out == nil {
		out = &Recorder{}
	}
	return New(id,
		name,
		out,
		&rules.Lineage{},
		&rules.Class{},
		rules.Abilities{})
}

func (p *Player) Send(msg interface{}) error {
	return p.out.Send(msg)
}

func (p *Player) TakeMeleeDamage(damage int64) bool {
	p.isDirty = true
	p.curHealth -= damage
	if p.curHealth <= 0 {
		return true
	}
	return false
}

func (p *Player) RestoreHealth(amount int64) {
	p.curHealth = min(p.curHealth+amount, p.maxHealth)
}

func (p *Player) IsDead() bool {
	return p.curHealth <= 0
}

func (p *Player) ResetDirtyFlag() {
	p.isDirty = false
	p.inventory.ResetDirtyFlag()
	p.slots.ResetDirtyFlag()
}
