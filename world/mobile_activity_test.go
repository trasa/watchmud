package world

import (
	"github.com/stretchr/testify/assert"
	"github.com/trasa/watchmud/direction"
	"github.com/trasa/watchmud/mobile"
	"testing"
	"time"
)

func Test_getNextDirectionOnPath_Simple(t *testing.T) {
	m := mobile.NewInstance("id",
		mobile.NewDefinition("id", "name", []string{}, "desc", "room desc", mobile.WanderingDefinition{
			CanWander:       true,
			CheckFrequency:  time.Minute * 1,
			CheckPercentage: 1.0,
			Style:           mobile.WANDER_FOLLOW_PATH,
			Path:            []string{"a", "b"},
		}))

	r := NewTestRoom("a")
	r.Up = NewTestRoom("b")
	r.Up.Down = r

	// a -> b
	dir, changeDirection, err := getNextDirectionOnPath(m, r)
	assert.NoError(t, err)
	assert.Equal(t, direction.UP, dir)
	assert.False(t, changeDirection)

	// b -> a
	dir, changeDirection, err = getNextDirectionOnPath(m, r.Up)
	assert.NoError(t, err)
	assert.Equal(t, direction.DOWN, dir)
	assert.True(t, changeDirection)
}

func Test_getNextDirectionOnPath_FullPath(t *testing.T) {
	m := mobile.NewInstance("id",
		mobile.NewDefinition("id", "name", []string{}, "desc", "room desc", mobile.WanderingDefinition{
			CanWander:       true,
			CheckFrequency:  time.Minute * 1,
			CheckPercentage: 1.0,
			Style:           mobile.WANDER_FOLLOW_PATH,
			Path:            []string{"a", "b", "c"},
		}))

	// a <-> b <-> c
	a := NewTestRoom("a")
	b := NewTestRoom("b")
	c := NewTestRoom("c")
	a.East = b
	b.West = a
	b.East = c
	c.West = b

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
