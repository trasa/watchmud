package world

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRoomExits_none(t *testing.T) {
	r := &Room{}
	exits := r.GetExits()
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

	exits := r.GetExits()
	assert.Equal(t, "neswud", exits)
}

func TestRoomExits_some(t *testing.T) {
	r := &Room{}
	r.North = r
	r.East = r
	r.Up = r

	exits := r.GetExits()
	assert.Equal(t, "neu", exits)
}
