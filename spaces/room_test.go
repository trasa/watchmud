package spaces

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/trasa/watchmud-message/direction"
)

func TestRoomExits_none(t *testing.T) {
	r := NewTestRoom("testing")
	exits := r.Exits()
	assert.Equal(t, "", exits)
}

func TestRoomExits_all(t *testing.T) {
	r := NewTestRoom("testing")
	r.ConnectRoom(direction.North, r)
	r.ConnectRoom(direction.South, r)
	r.ConnectRoom(direction.East, r)
	r.ConnectRoom(direction.West, r)
	r.ConnectRoom(direction.Up, r)
	r.ConnectRoom(direction.Down, r)

	exits := r.Exits()
	assert.Equal(t, "neswud", exits)
}

func TestRoomExits_some(t *testing.T) {
	r := NewTestRoom("test")
	r.ConnectRoom(direction.North, r)
	r.ConnectRoom(direction.East, r)
	r.ConnectRoom(direction.Up, r)

	exits := r.Exits()
	assert.Equal(t, "neu", exits)
}

func TestRoom_GetExitInfo(t *testing.T) {
	center := NewTestRoom("center")
	n := NewTestRoom("n")
	s := NewTestRoom("s")

	center.ConnectRoom(direction.North, n)
	n.ConnectRoom(direction.South, center)

	center.ConnectRoom(direction.South, s)
	s.ConnectRoom(direction.North, center)

	exitInfo := center.GetRoomExits(false)

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
	center.ConnectRoom(direction.North, n)
	// one choice
	dir = center.PickRandomDirection(false)
	assert.Equal(t, direction.North, dir)

	// two choices
	s := NewTestRoom("s")
	center.ConnectRoom(direction.South, s)

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

	center.ConnectRoom(direction.North, n)
	n.ConnectRoom(direction.South, center)

	center.ConnectRoom(direction.South, s)
	s.ConnectRoom(direction.North, center)

	result := center.GetRoomExits(true)
	assert.Equal(t, 1, len(result))
	assert.Equal(t, direction.North, result[0].Direction)
}
