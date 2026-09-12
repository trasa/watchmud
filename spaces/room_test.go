package spaces

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/trasa/watchmud/direction"
)

func TestRoomExits_none(t *testing.T) {
	r := NewTestRoom("testing")
	exits := r.ExitString()
	assert.Equal(t, "None!", exits)
}

func TestRoomExits_all(t *testing.T) {
	r := NewTestRoom("testing")
	r.Connect(direction.North, r)
	r.Connect(direction.South, r)
	r.Connect(direction.East, r)
	r.Connect(direction.West, r)
	r.Connect(direction.Up, r)
	r.Connect(direction.Down, r)

	exits := r.ExitString()
	assert.Equal(t, "North, East, South, West, Up, Down", exits)
}

func TestRoomExits_some(t *testing.T) {
	r := NewTestRoom("test")
	r.Connect(direction.North, r)
	r.Connect(direction.East, r)
	r.Connect(direction.Up, r)

	exits := r.ExitString()
	assert.Equal(t, "North, East, Up", exits)
}

func TestRoom_GetExitInfo(t *testing.T) {
	center := NewTestRoom("center")
	n := NewTestRoom("n")
	s := NewTestRoom("s")

	center.Connect(direction.North, n)
	n.Connect(direction.South, center)

	center.Connect(direction.South, s)
	s.Connect(direction.North, center)

	exitInfo := center.Exits(false)

	assert.Equal(t, 2, len(exitInfo))
	assert.Equal(t, direction.North, exitInfo[0].Direction)
	assert.Equal(t, direction.South, exitInfo[1].Direction)
}

func TestRoom_PickRandomDirection(t *testing.T) {
	center := NewTestRoom("center")
	// no rooms out
	dir := center.PickRandomDirection(false)
	assert.Equal(t, direction.None, dir)

	n := NewTestRoom("n")
	center.Connect(direction.North, n)
	// one choice
	dir = center.PickRandomDirection(false)
	assert.Equal(t, direction.North, dir)

	// two choices
	s := NewTestRoom("s")
	center.Connect(direction.South, s)

	dir = center.PickRandomDirection(false)
	if !(dir == direction.North || dir == direction.South) {
		t.Errorf("expected NORTH or SOUTH but found %d", dir)
	}
}

func TestRoom_LimitToZone(t *testing.T) {
	zone1 := &Zone{Id: "zone1"}
	zone2 := &Zone{Id: "zone2"}
	center := NewTestRoom("center")
	center.Zone = zone1

	n := NewTestRoom("n")
	n.Zone = zone1
	s := NewTestRoom("s")
	s.Zone = zone2

	center.Connect(direction.North, n)
	n.Connect(direction.South, center)

	center.Connect(direction.South, s)
	s.Connect(direction.North, center)

	result := center.Exits(true)
	assert.Equal(t, 1, len(result))
	assert.Equal(t, direction.North, result[0].Direction)
}
