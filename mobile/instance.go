package mobile

import (
	"errors"
	"fmt"
	"github.com/satori/go.uuid"
	"github.com/trasa/watchmud/combat"
	"math/rand"
	"time"
)

// The Mobile standing in front of you is an Instance of
// its definition. Mobiles of definition 'lizard' are all
// immune to poison, but this instance of 'lizard' is wearing
// a magic hat and has a sword in it's hand. (scary lizard)
type Instance struct {
	InstanceId        uuid.UUID
	Definition        *Definition
	LastWanderingTime time.Time // when was the last time this mob went wandering?
	WanderingForward  bool      // do you wander forward on the path or backwards?
	CurHealth         int64
}

func NewInstance(defn *Definition) *Instance {
	return &Instance{
		InstanceId:        uuid.NewV4(),
		Definition:        defn,
		LastWanderingTime: time.Now(),
		WanderingForward:  true, // by default
		CurHealth:         defn.MaxHealth,
	}
}

func (mob *Instance) Id() string {
	return mob.InstanceId.String()
}

func (mob *Instance) GetName() string {
	return mob.Definition.Name
}

func (mob *Instance) CanWander() bool {
	return mob.canWander(time.Now())
}

func (mob *Instance) canWander(now time.Time) bool {
	if !mob.Definition.Wandering.CanWander {
		return false
	}
	timeSince := now.Sub(mob.LastWanderingTime)
	return timeSince > mob.Definition.Wandering.CheckFrequency
}

func (mob *Instance) CheckWanderChance() bool {
	return mob.checkWanderChance(rand.New(rand.NewSource(time.Now().UnixNano())))
}

func (mob *Instance) checkWanderChance(r *rand.Rand) bool {
	chance := r.Float32()
	//log.Printf("mob '%s' chance of walking %f vs. %f", mob.Definition.Id, chance, mob.Definition.Wandering.CheckPercentage)
	return chance < mob.Definition.Wandering.CheckPercentage
}

// Determine where we are on the wandering path given the current room id.
// returns error if we're not wandering on a path
func (mob *Instance) GetIndexOnPath(currentRoom string) (int, error) {
	if len(mob.Definition.Wandering.Path) == 0 {
		return -1, errors.New("instance not defined to be on a path")
	}
	for i, s := range mob.Definition.Wandering.Path {
		if s == currentRoom {
			return i, nil
		}
	}
	return -1, errors.New(fmt.Sprintf("currentRoom '%s' not found in path '%s'", currentRoom, mob.Definition.Wandering.Path))
}

// Combatant
func (mob *Instance) TakeMeleeDamage(damage int64) (isDead bool) {
	mob.CurHealth = mob.CurHealth - damage
	return mob.CurHealth <= 0
}

// Combatant
func (mob *Instance) IsDead() bool {
	return mob.CurHealth <= 0
}

// Combatant
func (mob *Instance) CombatantType() combat.CombatantType {
	return combat.MobileCombatant
}

// Combatant
func (mob *Instance) CalculateMeleeRollModifiers() int {
	// TODO mob melee modifiers come from definition and ...?
	return 0
}

// Combatant
func (mob *Instance) ArmorClass() int {
	// base +/- equipment
	// TODO
	return mob.Definition.ArmorClass()
}

// Combatant
func (mob *Instance) HasResistanceTo(damageType combat.DamageType) bool {
	// resistance gained by spells or equipment?
	// TODO
	return mob.Definition.HasResistanceTo(damageType)
}

// Combatant
func (mob *Instance) IsVulnerableTo(damageType combat.DamageType) bool {
	// vulnerability gained by spells or equipment?
	// TODO
	return mob.Definition.IsVulnerableTo(damageType)
}

// Combatant
func (mob *Instance) WeaponDamageRoll() string {
	// TODO mobs cant carry anything yet?
	return "1d6"
}

// Combatant
func (mob *Instance) WeaponDamageType() combat.DamageType {
	// TODO mobs can't carry anything yet
	return combat.Piercing
}

func (mob *Instance) Send(msg interface{}) {
	// TODO do something with this notification
}

// Restore mob's health, mana, movement
func (mob *Instance) Restore() {
	mob.CurHealth = mob.Definition.MaxHealth
}
