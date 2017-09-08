package direction

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStringToDirection_Successful(t *testing.T) {
	assertDirection := func(dir Direction, s string) {
		parsed, err := StringToDirection(s)
		if err != nil {
			t.Error(err)
		} else {
			assert.Equal(t, dir, parsed)
		}
	}
	assertDirection(NORTH, "n")
	assertDirection(NORTH, "north")
	assertDirection(NORTH, "NoRtH")
	assertDirection(EAST, "e")
	assertDirection(WEST, "w")
	assertDirection(SOUTH, "s")
	assertDirection(UP, "U")
	assertDirection(UP, "up")
	assertDirection(DOWN, "D")
}

func TestAbbreviationToDirection_EmptyString(t *testing.T) {
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

func TestAbbreviationToString(t *testing.T) {
	doit := func(dir string, expected string) {
		str, err := AbbreviationToString(dir)
		assert.Equal(t, expected, str)
		assert.Nil(t, err)
	}

	doit("n", "North")
	doit("s", "South")
	doit("e", "East")
	doit("w", "West")
	doit("u", "Up")
	doit("d", "Down")

	doit("N", "North")
	doit("S", "South")
	doit("E", "East")
	doit("W", "West")
	doit("U", "Up")
	doit("D", "Down")
}

func TestAbbreviationToString_IsUnknown(t *testing.T) {
	_, err := AbbreviationToString("x")
	assert.NotNil(t, err)
}

func TestAbbreviationToString_IsTooBig(t *testing.T) {
	_, err := AbbreviationToString("asdlfkjasdlf")
	assert.NotNil(t, err)
}

func TestExitsToString(t *testing.T) {
	full, err := ExitsToFormattedString("nsewud")
	assert.Equal(t, "North, South, East, West, Up, Down", full)
	assert.Nil(t, err)
}

func TestExitsToString_Empty(t *testing.T) {
	s, err := ExitsToFormattedString("")
	assert.Equal(t, "None!", s)
	assert.Nil(t, err)
}

func TestExitsToString_BadInput(t *testing.T) {
	s, err := ExitsToFormattedString("a")
	assert.Equal(t, "", s)
	assert.NotNil(t, err)

	s, err = ExitsToFormattedString("dasdflkj")
	assert.Equal(t, "", s)
	assert.NotNil(t, err)
}

func TestDirectionToAbbreviation(t *testing.T) {
	doit := func(dir Direction, expected string) {
		str, err := DirectionToAbbreviation(dir)
		assert.Equal(t, expected, str)
		assert.Nil(t, err)
	}

	doit(NORTH, "n")
	doit(SOUTH, "s")
	doit(EAST, "e")
	doit(WEST, "w")
	doit(UP, "u")
	doit(DOWN, "d")
}

func TestDirectionToString(t *testing.T) {
	doit := func(dir Direction, expected string) {
		str, err := DirectionToString(dir)
		assert.Equal(t, expected, str)
		assert.Nil(t, err)
	}

	doit(NORTH, "North")
	doit(SOUTH, "South")
	doit(EAST, "East")
	doit(WEST, "West")
	doit(UP, "Up")
	doit(DOWN, "Down")
}
