package player

import (
	"github.com/trasa/watchmud/combat"
	"github.com/trasa/watchmud/object"
)

// see https://play.golang.org/p/zPLyr3ZOM0 (first attempt)
// then see https://play.golang.org/p/z5athD5fV3 (client is an interface, but now pointer woes)
//noinspection GoNameStartsWithPackageName
type Player interface {
	// send a message to the player
	Send(innerMessage interface{}) error
	// get the unique player_id of this player
	GetId() int64
	// get the unique player_name of the player
	GetName() string
	// get the player's race id
	GetRaceId() int32
	// get the player's class id
	GetClassId() int32
	// get the player's inventory
	Inventory() *Inventory
	// get the player's slots (what they are wearing, where)
	Slots() *object.Slots
	// get current health
	GetCurrentHealth() int64
	// get max health
	GetMaxHealth() int64
	// apply damage to the player and see if they're dead
	TakeMeleeDamage(damage int64) (isDead bool)
	// are they dead?
	IsDead() bool
	// what type of combatant is this?
	CombatantType() combat.CombatantType
	// if true, we need to write player data back to database
	IsDirty() bool
	// set dirty flag to false (so we won't write back to database)
	ResetDirtyFlag()
	// where are you?
	Location() *Location
	// your ability scores
	Abilities() *Abilities
	// Figure out + and - to an attack roll
	CalculateMeleeRollModifiers() int
	// calculate the player's AC based on intrinsics and equipment
	ArmorClass() int
	// resistance to damage?
	HasResistanceTo(damageType combat.DamageType) bool
	// vulnerable to damage?
	IsVulnerableTo(damageType combat.DamageType) bool
	// string of the dice roll for wielded weapon
	WeaponDamageRoll() string
	// damage type for wielded weapon
	WeaponDamageType() combat.DamageType
	// restore health, mana, movement, ...
	Restore()
}
