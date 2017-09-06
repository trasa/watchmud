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

func TestAbbreviationToString(t *testing.T) {
	doit := func(dir string, expected string) {
		str, err := abbreviationToString(dir)
		assert.Equal(t, expected, str)
		assert.Nil(t, err)
	}

	doit("n", "North")
	doit("s", "South")
	doit("e", "East")
	doit("w", "West")
	doit("u", "Up")
	doit("d", "Down")
}

func TestAbbreviationIsUnknown(t *testing.T) {
	_, err := abbreviationToString("x")
	assert.NotNil(t, err)
}

func TestAbbreviationIsTooBig(t *testing.T) {
	_, err := abbreviationToString("asdlfkjasdlf")
	assert.NotNil(t, err)
}

func TestExitsToString(t *testing.T) {
	full, err := ExitsToString("nsewud")
	assert.Equal(t, "North, South, East, West, Up, Down", full)
	assert.Nil(t, err)
}

func TestExitsToStringEmpty(t *testing.T) {
	s, err := ExitsToString("")
	assert.Equal(t, "None!", s)
	assert.Nil(t, err)
}

func TestExitsToStringBadInput(t *testing.T) {
	s, err := ExitsToString("a")
	assert.Equal(t, "", s)
	assert.NotNil(t, err)

	s, err = ExitsToString("dasdflkj")
	assert.Equal(t, "", s)
	assert.NotNil(t, err)
}
