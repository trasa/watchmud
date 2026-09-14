package wandering

import "time"

// Definition describes how the mob wanders around the world.
type Definition struct {
	CanWander       bool
	Style           Style         // how do you wander?
	CheckFrequency  time.Duration // how long between wandering?
	CheckPercentage float32       // % chance of moving on each test
	Path            []string
}

type Style int

const (
	None       Style = iota // you don't wander
	Random                  // wander within the zone randomly
	FollowPath              // wander a prescribed path which could cross zones
)
