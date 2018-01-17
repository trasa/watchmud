package spaces

import (
	"github.com/stretchr/testify/assert"
	"github.com/trasa/watchmud/direction"
	"testing"
)

func TestRoomExits_none(t *testing.T) {
	r := NewTestRoom("testing")
	exits := r.GetExitString()
	assert.Equal(t, "", exits)
}

func TestRoomExits_all(t *testing.T) {
	r := NewTestRoom("testing")
	r.Set(direction.NORTH, r)
	r.Set(direction.SOUTH, r)
	r.Set(direction.EAST, r)
	r.Set(direction.WEST, r)
	r.Set(direction.UP, r)
	r.Set(direction.DOWN, r)

	exits := r.GetExitString()
	assert.Equal(t, "neswud", exits)
}

func TestRoomExits_some(t *testing.T) {
	r := NewTestRoom("test")
	r.Set(direction.NORTH, r)
	r.Set(direction.EAST, r)
	r.Set(direction.UP, r)

	exits := r.GetExitString()
	assert.Equal(t, "neu", exits)
}

func TestRoom_GetExitInfo(t *testing.T) {
	center := NewTestRoom("center")
	n := NewTestRoom("n")
	s := NewTestRoom("s")

	center.Set(direction.NORTH, n)
	n.Set(direction.SOUTH, center)

	center.Set(direction.SOUTH, s)
	s.Set(direction.NORTH, center)

	exitInfo := center.GetRoomExits(false)

	assert.Equal(t, 2, len(exitInfo))
	assert.Equal(t, direction.NORTH, exitInfo[0].Direction)
	assert.Equal(t, direction.SOUTH, exitInfo[1].Direction)
}

func TestRoom_PickRandomDirection(t *testing.T) {
	center := NewTestRoom("center")
	// no rooms out
	dir := center.PickRandomDirection(false)
	assert.Equal(t, direction.NONE, dir)

	n := NewTestRoom("n")
	center.Set(direction.NORTH, n)
	// one choice
	dir = center.PickRandomDirection(false)
	assert.Equal(t, direction.NORTH, dir)

	// two choices
	s := NewTestRoom("s")
	center.Set(direction.SOUTH, s)

	dir = center.PickRandomDirection(false)
	if !(dir == direction.NORTH || dir == direction.SOUTH) {
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

	center.Set(direction.NORTH, n)
	n.Set(direction.SOUTH, center)

	center.Set(direction.SOUTH, s)
	s.Set(direction.NORTH, center)

	result := center.GetRoomExits(true)
	assert.Equal(t, 1, len(result))
	assert.Equal(t, direction.NORTH, result[0].Direction)
}
