package direction

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringToDirection_Successful(t *testing.T) {
	assertDirection := func(dir Direction, s string) {
		parsed, err := Parse(s)
		if err != nil {
			t.Error(err)
		} else {
			assert.Equal(t, dir, parsed)
		}
	}
	assertDirection(North, "n")
	assertDirection(North, "north")
	assertDirection(North, "NoRtH")
	assertDirection(East, "e")
	assertDirection(West, "w")
	assertDirection(South, "s")
	assertDirection(Up, "U")
	assertDirection(Up, "up")
	assertDirection(Down, "D")
}

func TestAbbreviationToDirection_EmptyString(t *testing.T) {
	_, err := Parse("")
	if err == nil {
		t.Error("Expected error")
	}
}

func TestStringToDirection_Unknown(t *testing.T) {
	_, err := Parse("asdf")
	if err == nil {
		t.Error("expected error")
	}
}

func TestAbbreviationToString(t *testing.T) {
	doit := func(dir string, expected Direction) {
		d, err := ParsePrefix(dir)
		assert.Equal(t, expected, d)
		assert.Nil(t, err)
	}

	doit("n", North)
	doit("s", South)
	doit("e", East)
	doit("w", West)
	doit("u", Up)
	doit("d", Down)

	doit("N", North)
	doit("S", South)
	doit("E", East)
	doit("W", West)
	doit("U", Up)
	doit("D", Down)
}

func TestAbbreviationToString_IsUnknown(t *testing.T) {
	_, err := ParsePrefix("x")
	assert.NotNil(t, err)
}

func TestAbbreviationToString_IsTooBig(t *testing.T) {
	_, err := ParsePrefix("asdlfkjasdlf")
	assert.NotNil(t, err)
}

func TestDirectionToAbbreviation(t *testing.T) {
	doit := func(dir Direction, expected string) {
		str := dir.Abbrev()
		assert.Equal(t, expected, str)
	}

	doit(North, "n")
	doit(South, "s")
	doit(East, "e")
	doit(West, "w")
	doit(Up, "u")
	doit(Down, "d")
}

func TestDirectionToString(t *testing.T) {
	doit := func(dir Direction, expected string) {
		str := dir.String()
		assert.Equal(t, expected, str)
	}

	doit(North, "North")
	doit(South, "South")
	doit(East, "East")
	doit(West, "West")
	doit(Up, "Up")
	doit(Down, "Down")
}
