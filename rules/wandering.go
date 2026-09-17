package rules

import (
	"errors"
	"fmt"
	"slices"
	"time"
)

// WanderDefinition describes how the mob wanders around the world.
type WanderDefinition struct {
	CanWander       bool
	Style           WanderStyle   // how do you wander?
	CheckFrequency  time.Duration // how long between wandering?
	CheckPercentage float32       // % chance of moving on each test
	Path            []string
}

type WanderStyle string

var ErrUnknownWanderStyle = errors.New("unknown armor type")

const (
	WanderNone       WanderStyle = ""           // you don't wander
	WanderRandom     WanderStyle = "random"     // wander within the zone randomly
	WanderFollowPath WanderStyle = "followPath" // wander a prescribed path which could cross zones
)

var wanderStyles = []WanderStyle{
	WanderNone,
	WanderRandom,
	WanderFollowPath,
}

func (t *WanderStyle) UnmarshalText(b []byte) error {
	s, err := ParseWanderStyle(string(b))
	if err != nil {
		return err
	}
	*t = s
	return nil
}

func ParseWanderStyle(s string) (WanderStyle, error) {
	ws := WanderStyle(s)
	if !slices.Contains(wanderStyles, ws) {
		return WanderNone, fmt.Errorf("%w: %q", ErrUnknownWanderStyle, s)
	}
	return ws, nil
}
