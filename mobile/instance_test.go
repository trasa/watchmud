package mobile

import (
	"github.com/stretchr/testify/assert"
	"log"
	"math/rand"
	"testing"
	"time"
)

func TestInstance_CanWander(t *testing.T) {
	noWalk := NewInstance(NewDefinition("id", "nowalk", "testZone", []string{}, "", "",
		25,
		WanderingDefinition{
			CanWander: false,
		}, 10))
	assert.False(t, noWalk.CanWander())

	walker := NewInstance(NewDefinition("", "", "", []string{}, "", "",
		25,
		WanderingDefinition{
			CanWander:       true,
			Style:           WANDER_RANDOM,
			CheckFrequency:  time.Minute * 1,
			CheckPercentage: 1.0,
		}, 10))
	log.Printf("wandering %s", walker.LastWanderingTime)
	assert.False(t, walker.CanWander()) // because it hasn't been 1 minute yet

	now := time.Now()
	walker.LastWanderingTime = now
	assert.False(t, walker.canWander(now))                     // can't wander right now
	assert.False(t, walker.canWander(now.Add(time.Second*30))) // 30 seconds from now
	assert.False(t, walker.canWander(now.Add(time.Second*60))) // 1 minute
	assert.True(t, walker.canWander(now.Add(time.Second*61)))  // 1 minute
}

func TestInstance_CheckWanderChance_AlwaysFails(t *testing.T) {
	r := rand.New(rand.NewSource(1))

	noChance := NewInstance(NewDefinition("", "", "", []string{}, "", "",
		25,
		WanderingDefinition{
			CanWander:       true,
			Style:           WANDER_RANDOM,
			CheckFrequency:  time.Minute * 1,
			CheckPercentage: 0.0, // <-- 0% chance
		}, 10))
	for i := 0; i < 10; i++ {
		assert.False(t, noChance.checkWanderChance(r)) // always fails
	}
}

func TestInstance_CheckWanderChance_AlwaysSucceeds(t *testing.T) {
	r := rand.New(rand.NewSource(1))
	fiftyFifty := NewInstance(NewDefinition("", "", "", []string{}, "", "",
		25,
		WanderingDefinition{
			CanWander:       true,
			Style:           WANDER_RANDOM,
			CheckFrequency:  time.Minute * 1,
			CheckPercentage: 1.0, // <-- 100% chance
		}, 10))
	for i := 0; i < 10; i++ {
		assert.True(t, fiftyFifty.checkWanderChance(r))
	}
}

func TestInstance_CheckWanderChance_FiftyFifty(t *testing.T) {
	r := rand.New(rand.NewSource(1))
	fiftyFifty := NewInstance(NewDefinition("", "", "", []string{}, "", "",
		25,
		WanderingDefinition{
			CanWander:       true,
			Style:           WANDER_RANDOM,
			CheckFrequency:  time.Minute * 1,
			CheckPercentage: 0.50, // <-- 50% chance
		}, 10))
	success := false
	for i := 0; i < 100; i++ {
		if fiftyFifty.checkWanderChance(r) {
			log.Printf("Success after %d attempts", i+1)
			success = true
			break
		}
	}
	assert.True(t, success)
}

func TestInstance_GetIndexOnPath(t *testing.T) {
	m := NewInstance(
		NewDefinition("id", "name", "", []string{}, "desc", "room desc",
			25,
			WanderingDefinition{
				CanWander:       true,
				CheckFrequency:  time.Minute * 1,
				CheckPercentage: 1.0,
				Style:           WANDER_FOLLOW_PATH,
				Path:            []string{"a", "b"},
			}, 10))
	idx, err := m.GetIndexOnPath("a")
	assert.NoError(t, err)
	assert.Equal(t, 0, idx)

	idx, err = m.GetIndexOnPath("b")
	assert.NoError(t, err)
	assert.Equal(t, 1, idx)

	idx, err = m.GetIndexOnPath("foo")
	assert.Error(t, err)
	assert.Equal(t, -1, idx)
}

func TestInstance_GetIndexOnPath_NoPath(t *testing.T) {
	m := NewInstance(
		NewDefinition("id", "name", "", []string{}, "desc", "room desc",
			25,
			WanderingDefinition{}, 10))
	idx, err := m.GetIndexOnPath("foo")
	assert.Error(t, err)
	assert.Equal(t, -1, idx)
}
