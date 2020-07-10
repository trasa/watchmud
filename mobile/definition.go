package mobile

import (
	"github.com/trasa/watchmud/combat"
	"time"
)

// Defines what it means to be a mob.
type Definition struct {
	Id                string
	Aliases           []string
	Name              string
	ShortDescription  string
	DescriptionInRoom string // description when in a room "A giant lizard is here."
	ZoneId            string
	Wandering         WanderingDefinition
	MaxHealth         int64
	flags             map[string]bool
	ac                int
}

func NewDefinition(definitionId string,
	name string,
	zoneId string,
	aliases []string,
	shortDescription,
	descriptionInRoom string,
	maxHealth int64,
	wandering WanderingDefinition) *Definition {
	d := &Definition{
		Id:                definitionId,
		Name:              name,
		Aliases:           aliases,
		ShortDescription:  shortDescription,
		DescriptionInRoom: descriptionInRoom,
		Wandering:         wandering,
		ZoneId:            zoneId,
		flags:             make(map[string]bool),
		MaxHealth:         maxHealth,
	}
	return d
}

func (d *Definition) HasAlias(target string) bool {
	for _, a := range d.Aliases {
		if target == a {
			return true
		}
	}
	return false
}

func (d *Definition) SetFlags(flags []string) {
	if flags != nil {
		for _, s := range flags {
			d.SetFlag(s)
		}
	}
}

func (d *Definition) SetFlag(flag string) {
	d.flags[flag] = true
}

func (d *Definition) HasFlag(flag string) bool {
	return d.flags[flag]
}

func (d *Definition) GetFlags() (result []string) {
	for k, v := range d.flags {
		if v {
			result = append(result, k)
		}
	}
	return
}

func (d *Definition) ArmorClass() int {
	// base armor class before modifiers
	return d.ac
}

func (d *Definition) HasResistanceTo(damageType combat.DamageType) bool {
	// TODO resistances
	return false
}

func (d *Definition) IsVulnerableTo(damageType combat.DamageType) bool {
	// TODO vulnerability
	return false
}

// Things to do with how mobs wander around
type WanderingDefinition struct {
	CanWander       bool
	Style           WanderingStyle // how do you wander?
	CheckFrequency  time.Duration  // how long between wandering?
	CheckPercentage float32        // % chance of moving on each test
	Path            []string
}

type WanderingStyle int

const (
	WANDER_NONE        WanderingStyle = iota // you don't wander
	WANDER_RANDOM                            // wander within the zone randomly
	WANDER_FOLLOW_PATH                       // wander a prescribed path which could cross zones
)
