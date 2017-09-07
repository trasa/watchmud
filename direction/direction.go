package direction

import (
	"errors"
	"fmt"
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

// turn an abbreviation like "n" into
// a direction string like "North"
func AbbreviationToString(abbrev string) (string, error) {
	if len(abbrev) != 1 {
		return "", errors.New("Expected abbrev to be 1 character")
	}
	switch abbrev {
	case "n":
		return "North", nil
	case "s":
		return "South", nil
	case "e":
		return "East", nil
	case "w":
		return "West", nil
	case "u":
		return "Up", nil
	case "d":
		return "Down", nil
	default:
		return "", errors.New(fmt.Sprintf("Unhandled Direction Abbreviation: %s", abbrev))
	}
}

// Given an 'exits' string like "nsed", return a formatted
// string for the player with something like
// "North, South, East, Down"
func ExitsToString(exits string) (result string, err error) {
	if len(exits) == 0 {
		result = "None!"
	} else {
		for _, dir := range strings.Split(exits, "") {
			var str string
			if str, err = AbbreviationToString(dir); err != nil {
				result = ""
				return
			}
			result += str + ", "
		}
		result = result[:len(result)-2]
	}
	return result, err
}
