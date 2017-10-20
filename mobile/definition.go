package mobile

import "time"

func NewDefinition(definitionId string,
	name string,
	aliases []string,
	shortDescription,
	descriptionInRoom string,
	wandering WanderingDefinition) *Definition {
	d := &Definition{
		DefinitionId:      definitionId,
		Name:              name,
		Aliases:           aliases,
		ShortDescription:  shortDescription,
		DescriptionInRoom: descriptionInRoom,
		Wandering:         wandering,
	}
	return d
}

// Defines what it means to be a mob.
type Definition struct {
	DefinitionId      string
	Aliases           []string
	Name              string
	ShortDescription  string
	DescriptionInRoom string // description when in a room "A giant lizard is here."
	Wandering         WanderingDefinition
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
