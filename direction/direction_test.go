package direction

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStringToDirection_Successful(t *testing.T) {
	assertDirection(t, NORTH, "n")
	assertDirection(t, NORTH, "north")
	assertDirection(t, NORTH, "NoRtH")
	assertDirection(t, EAST, "e")
	assertDirection(t, WEST, "w")
	assertDirection(t, SOUTH, "s")
	assertDirection(t, UP, "U")
	assertDirection(t, DOWN, "D")
}

func TestStringToDirection_EmptyString(t *testing.T) {
	_, err := StringToDirection("")
	if err == nil {
		t.Error("Expected error")
	}
}

func TestStringToDirection_Unknown(t *testing.T) {
	_, err := StringToDirection("asdf")
	if err == nil {
		t.Error("expected error")
	}
}

func assertDirection(t *testing.T, dir Direction, s string) {
	parsed, err := StringToDirection(s)
	if err != nil {
		t.Error(err)
	} else {
		assert.Equal(t, dir, parsed)
	}
}
