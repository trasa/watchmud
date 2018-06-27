package mobile

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

// The Mobile standing in front of you is an Instance of
// its definition. Mobiles of definition 'lizard' are all
// immune to poison, but this instance of 'lizard' is wearing
// a magic hat and has a sword in it's hand. (scary lizard)
type Instance struct {
	Definition        *Definition
	LastWanderingTime time.Time // when was the last time this mob went wandering?
	WanderingForward  bool      // do you wander forward on the path or backwards?
	CurHealth         int
}

func NewInstance(defn *Definition) *Instance {
	return &Instance{
		Definition:        defn,
		LastWanderingTime: time.Now(),
		WanderingForward:  true, // by default
		CurHealth:         defn.MaxHealth,
	}
}

func (mob *Instance) CanWander() bool {
	return mob.canWander(time.Now())
}

func (mob *Instance) canWander(now time.Time) bool {
	if !mob.Definition.Wandering.CanWander {
		return false
	}
	// TODO don't wander off if you're in a fight
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
