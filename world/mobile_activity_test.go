package world

import (
	"github.com/stretchr/testify/assert"
	"github.com/trasa/watchmud-message/direction"
	"github.com/trasa/watchmud/mobile"
	"github.com/trasa/watchmud/spaces"
	"testing"
	"time"
)

func Test_getNextDirectionOnPath_Simple(t *testing.T) {
	m := mobile.NewInstance(
		mobile.NewDefinition("id", "name", "", []string{}, "desc", "room desc", mobile.WanderingDefinition{
			CanWander:       true,
			CheckFrequency:  time.Minute * 1,
			CheckPercentage: 1.0,
			Style:           mobile.WANDER_FOLLOW_PATH,
			Path:            []string{"a", "b"},
		}))

	r := spaces.NewTestRoom("a")
	r.Set(direction.UP, spaces.NewTestRoom("b"))
	r.Get(direction.UP).Set(direction.DOWN, r)

	// a -> b
	dir, changeDirection, err := getNextDirectionOnPath(m, r)
	assert.NoError(t, err)
	assert.Equal(t, direction.UP, dir)
	assert.False(t, changeDirection)

	// b -> a
	dir, changeDirection, err = getNextDirectionOnPath(m, r.Get(direction.UP))
	assert.NoError(t, err)
	assert.Equal(t, direction.DOWN, dir)
	assert.True(t, changeDirection)
}

func Test_getNextDirectionOnPath_FullPath(t *testing.T) {
	m := mobile.NewInstance(
		mobile.NewDefinition("id", "name", "", []string{}, "desc", "room desc", mobile.WanderingDefinition{
			CanWander:       true,
			CheckFrequency:  time.Minute * 1,
			CheckPercentage: 1.0,
			Style:           mobile.WANDER_FOLLOW_PATH,
			Path:            []string{"a", "b", "c"},
		}))

	// a <-> b <-> c
	a := spaces.NewTestRoom("a")
	b := spaces.NewTestRoom("b")
	c := spaces.NewTestRoom("c")
	a.Set(direction.EAST, b)
	b.Set(direction.WEST, a)
	b.Set(direction.EAST, c)
	c.Set(direction.WEST, b)

	// a -> b
	dir, changeDirection, err := getNextDirectionOnPath(m, a)
	assert.NoError(t, err)
	assert.Equal(t, direction.EAST, dir)
	assert.False(t, changeDirection)

	// b -> c
	dir, changeDirection, err = getNextDirectionOnPath(m, b)
	assert.NoError(t, err)
	assert.Equal(t, direction.EAST, dir)
	assert.False(t, changeDirection)

	// c -> b
	dir, changeDirection, err = getNextDirectionOnPath(m, c)
	assert.NoError(t, err)
	assert.Equal(t, direction.WEST, dir)
	assert.True(t, changeDirection)

	// b -> a
	// mob needs to be walking back for this to work
	m.WanderingForward = false
	dir, changeDirection, err = getNextDirectionOnPath(m, b)
	assert.NoError(t, err)
	assert.Equal(t, direction.WEST, dir)
	assert.False(t, changeDirection) // since we're already walking backwards
}
