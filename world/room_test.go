package world

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func NewTestRoom(name string) *Room {
	return NewRoom(nil, name, name, "")
}

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

	exitInfo := center.GetExitInfo()

	assert.Equal(t, 2, len(exitInfo))
	assert.Equal(t, "n", exitInfo["n"])
	assert.Equal(t, "s", exitInfo["s"])
}
