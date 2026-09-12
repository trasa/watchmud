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
	id   uuid.UUID
	name string
	out  Sender // was: client.Client, via ClientPlayer

	// Lineage is cosmetic and nothing reads it but the renderer. There is no
	// Class beside it any more and no Role in its place: a role is read off
	// the equipped slots every time it is asked for, never stored.
	Lineage   *rules.Lineage
	inventory *Inventory
	slots     *Slots
	curHealth int64
	maxHealth int64
	location  Location
}

func New(id uuid.UUID,
	name string,
	out Sender,
	lineage *rules.Lineage,
) *Player {
	return &Player{
		id:        id,
		name:      name,
		out:       out,
		Lineage:   lineage,
		inventory: NewInventory(),
		slots:     NewSlots(),
		curHealth: 100, // TODO need a default here,
		maxHealth: 100,
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

// RoleWeights totals what the player's equipped gear contributes to each
// role. Resolving that to a rules.Role is the caller's job -- see
// rules.Catalog.RoleFor.
func (p *Player) RoleWeights() map[string]int {
	return p.slots.RoleWeights()
}

// LineageName is the player's lineage for display. Cosmetic, and empty if
// they somehow have none.
func (p *Player) LineageName() string {
	if p.Lineage == nil {
		return ""
	}
	return p.Lineage.Name
}

func (p *Player) CurrentHealth() int64 {
	return p.curHealth
}

func (p *Player) MaxHealth() int64 {
	return p.maxHealth
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
