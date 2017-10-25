package spaces

import (
	"github.com/stretchr/testify/assert"
	"github.com/trasa/watchmud/direction"
	"testing"
)

func TestRoomExits_none(t *testing.T) {
	r := &Room{}
	exits := r.GetExitString()
	assert.Equal(t, "", exits)
}

func TestRoomExits_all(t *testing.T) {
	r := &Room{}
	r.North = r
	r.South = r
	r.East = r
	r.West = r
	r.Up = r
	r.Down = r

	exits := r.GetExitString()
	assert.Equal(t, "neswud", exits)
}

func TestRoomExits_some(t *testing.T) {
	r := &Room{}
	r.North = r
	r.East = r
	r.Up = r

	exits := r.GetExitString()
	assert.Equal(t, "neu", exits)
}

func TestRoom_GetExitInfo(t *testing.T) {
	center := NewTestRoom("center")
	n := NewTestRoom("n")
	s := NewTestRoom("s")

	center.North = n
	n.South = center

	center.South = s
	s.North = center

	exitInfo := center.GetExitInfo(false)

	assert.Equal(t, 2, len(exitInfo))
	assert.Equal(t, "n", exitInfo[direction.NORTH])
	assert.Equal(t, "s", exitInfo[direction.SOUTH])
}

func TestRoom_PickRandomDirection(t *testing.T) {
	center := NewTestRoom("center")
	// no rooms out
	dir := center.PickRandomDirection(false)
	assert.Equal(t, direction.NONE, dir)

	n := NewTestRoom("n")
	center.North = n
	// one choice
	dir = center.PickRandomDirection(false)
	assert.Equal(t, direction.NORTH, dir)

	// two choices
	s := NewTestRoom("s")
	center.South = s

	dir = center.PickRandomDirection(false)
	if !(dir == direction.NORTH || dir == direction.SOUTH) {
		t.Errorf("expected NORTH or SOUTH but found %s", dir)
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

	center.North = n
	n.South = center

	center.South = s
	s.North = center

	result := center.GetExitInfo(true)
	assert.Equal(t, 1, len(result))
	assert.Equal(t, "n", result[direction.NORTH])
	assert.Equal(t, "", result[direction.SOUTH])
}
