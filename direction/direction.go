package direction

import (
	"errors"
	"strings"
)

type Direction int

const (
	NORTH Direction = iota
	EAST
	SOUTH
	WEST
	UP
	DOWN
)

func StringToDirection(s string) (dir Direction, err error) {
	if len(s) == 0 {
		err = errors.New("invalid direction string")
		return
	}
	switch strings.ToLower(s[:1]) {
	case "n":
		dir = NORTH
	case "e":
		dir = EAST
	case "w":
		dir = WEST
	case "s":
		dir = SOUTH
	case "u":
		dir = UP
	case "d":
		dir = DOWN
	default:
		err = errors.New("Unknown direction: " + strings.ToLower(s[:1]))
	}
	return
}

func DirectionToString(dir Direction) (str string, err error) {
	switch dir {
	case NORTH:
		str = "n"
	case SOUTH:
		str = "s"
	case EAST:
		str = "e"
	case WEST:
		str = "w"
	case UP:
		str = "u"
	case DOWN:
		str = "d"
	default:
		err = errors.New("Unknown direction")
	}
	return
}
