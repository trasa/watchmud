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
	Id        int64
	Name      string
	Client    client.Client
	inventory *player.Inventory
	slots     *object.Slots
	curHealth int64
	maxHealth int64
	race      int32
	class     int32
	dirty     bool
	location  *player.Location
	abilities *player.Abilities
}

// Create a ClientPlayer connected to a new TestClient
// (for testing)
func NewTestClientPlayer(name string) (p *ClientPlayer, cli *client.TestClient) {
	p = NewClientPlayer(name, 0, 0, player.NewAbilities(15, 14, 13, 12, 10, 9), nil)
	cli = client.NewTestClient(p)
	p.Client = cli
	return
}

// Create a new player and set it up to work with this client
func NewClientPlayer(name string, race int32, class int32, abilities *player.Abilities, client client.Client) *ClientPlayer {
	p := ClientPlayer{
		Name:      name,
		Client:    client, // address of interface
		inventory: player.NewInventory(),
		slots:     object.NewSlots(),
		curHealth: 100, // TODO need a default health here
		maxHealth: 100,
		race:      race,
		class:     class,
		abilities: abilities,
	}
	return &p
}

// Load player information into this struct without flagging anything as dirty
// Does not load slot information as that has to happen after inventory is loaded
func NewClientPlayerFromPlayerData(name string, pd *db.PlayerData, client client.Client) *ClientPlayer {
	p := NewClientPlayer(name, pd.Race, pd.Class, player.NewAbilities(pd.Strength, pd.Dexterity, pd.Constitution, pd.Intelligence, pd.Wisdom, pd.Charisma), client)
	p.Id = pd.Id
	p.curHealth = pd.CurHealth
	p.maxHealth = pd.MaxHealth
	p.location = player.NewLocation(pd.LastZoneId, pd.LastRoomId)
	return p
}

// Load player inventory into this struct without flagging anything as dirty
func (p *ClientPlayer) LoadInventory(instance *object.Instance) {
	p.inventory.Load(instance)
}

func (p *ClientPlayer) GetId() int64 {
	return p.Id
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

func (p *ClientPlayer) GetRaceId() int32 {
	return p.race
}

func (p *ClientPlayer) GetClassId() int32 {
	return p.class
}

func (p *ClientPlayer) Inventory() *player.Inventory {
	return p.inventory
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

func (p *ClientPlayer) CalculateMeleeRollModifiers() int {
	// TODO figure out modifiers based on stats and weapons and things
	return 0
}

func (p *ClientPlayer) ArmorClass() int {
	// TODO figure out modifiers based on stats and weapons and things
	return 10
}

func (p *ClientPlayer) HasResistanceTo(damageType combat.DamageType) bool {
	return false // TODO
}

func (p *ClientPlayer) IsVulnerableTo(damageType combat.DamageType) bool {
	return false // TODO
}

func (p *ClientPlayer) WeaponDamageRoll() string {
	// TODO get this from wielded weapon
	return "3d6"
}

func (p *ClientPlayer) WeaponDamageType() combat.DamageType {
	// TODO get this from wielded weapon
	return combat.Piercing
}

func (p *ClientPlayer) ResetDirtyFlag() {
	p.dirty = false
	p.inventory.ResetDirtyFlag()
	p.slots.ResetDirtyFlag()
}

func (p *ClientPlayer) IsDirty() bool {
	return p.dirty || p.inventory.IsDirty() || p.slots.IsDirty()
}

func (p *ClientPlayer) Location() *player.Location {
	return p.location
}

func (p *ClientPlayer) Abilities() *player.Abilities {
	return p.abilities
}
