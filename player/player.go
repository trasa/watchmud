package player

import (
	"uuid"

	"github.com/trasa/watchmud/combat"
	"github.com/trasa/watchmud/rules"
)

// Sender is anything that can deliver a message to this player's connection
type Sender interface {
	Send(msg any)
}

type Player struct {
	id        uuid.UUID
	name      string
	out       Sender // was: client.Client, via ClientPlayer
	Lineage   *rules.Lineage
	Class     *rules.Class
	inventory *Inventory
	slots     *Slots
	curHealth int64
	maxHealth int64
	location  Location
	abilities rules.Abilities
}

func New(id uuid.UUID,
	name string,
	out Sender,
	lineage *rules.Lineage,
	class *rules.Class,
	abilities rules.Abilities,
) *Player {
	return &Player{
		id:        id,
		name:      name,
		out:       out,
		Lineage:   lineage,
		Class:     class,
		inventory: NewInventory(),
		slots:     NewSlots(),
		curHealth: 100, // TODO need a default here,
		maxHealth: 100,
		abilities: abilities,
	}
}

func (p *Player) Id() uuid.UUID {
	return p.id
}

func (p *Player) Name() string {
	return p.name
}

// Inventory returns the inventory
func (p *Player) Inventory() *Inventory {
	// TODO is this needed? Should p.Inventory become visible?
	// is needing this call indicating a problem?
	return p.inventory
}

// Slots returns the inventory
func (p *Player) Slots() *Slots {
	// TODO is this needed? Should p.Inventory become visible?
	// is needing this call indicating a problem?
	return p.slots
}

func (p *Player) Send(msg any) {
	p.out.Send(msg)
}

func (p *Player) RestoreHealth(amount int64) {
	p.curHealth = min(p.curHealth+amount, p.maxHealth)
}

func (p *Player) Dead() bool {
	return p.curHealth <= 0
}

func (p *Player) ArmorClass() int {
	return p.slots.ArmorClass()
}

func (p *Player) CalculateMeleeRollModifiers() int {
	return 0
}

func (p *Player) HasResistanceTo(damageType combat.DamageType) bool {
	// TODO
	return false
}

func (p *Player) TakeMeleeDamage(damage int64) bool {
	p.curHealth -= damage
	if p.curHealth <= 0 {
		return true
	}
	return false
}

func (p *Player) IsVulnerableTo(damageType combat.DamageType) bool {
	// TODO
	return false
}

func (p *Player) WeaponDamageRoll() string {
	// TODO
	return "1d6"
}

func (p *Player) WeaponDamageType() combat.DamageType {
	// TODO
	return combat.Piercing
}
