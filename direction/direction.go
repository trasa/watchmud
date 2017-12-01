package direction

import (
	"errors"
	"strings"
)

//noinspection GoNameStartsWithPackageName
type Direction int32

const (
	NONE  Direction = iota // 0
	NORTH                  // 1
	EAST                   // 2
	SOUTH                  // 3
	WEST                   // 4
	UP                     // 5
	DOWN                   // 6
)

var oneLetterToDirection = map[string]Direction{
	"n": NORTH,
	"s": SOUTH,
	"e": EAST,
	"w": WEST,
	"u": UP,
	"d": DOWN,
}

//noinspection GoNameStartsWithPackageName
var directionToOneLetter = map[Direction]string{
	NORTH: "n",
	SOUTH: "s",
	EAST:  "e",
	WEST:  "w",
	UP:    "u",
	DOWN:  "d",
}

//noinspection GoNameStartsWithPackageName
var directionToWord = map[Direction]string{
	NONE:  "None!",
	NORTH: "North",
	SOUTH: "South",
	EAST:  "East",
	WEST:  "West",
	UP:    "Up",
	DOWN:  "Down",
}

// take an abbreviation such as "n" and turn it into a direction int const like Direction.NORTH
func StringToDirection(s string) (dir Direction, err error) {
	if len(s) == 0 {
		err = errors.New("invalid direction string")
		return
	}
	lookup := strings.ToLower(s[:1])
	dir = oneLetterToDirection[lookup]
	if dir == NONE {
		err = errors.New("Unknown direction: " + lookup)
	}
	return
}

// take a value such as Direction.NORTH and turn it into an abbreviation such as "n"
//noinspection GoNameStartsWithPackageName
func DirectionToAbbreviation(dir Direction) (str string, err error) {
	str = directionToOneLetter[dir]
	if str == "" {
		err = errors.New("Unknown direction")
	}
	return
}

// turn an abbreviation like "n" into a string like "North"
func AbbreviationToString(abbrev string) (str string, err error) {
	if len(abbrev) != 1 {
		return "", errors.New("Expected abbrev to be 1 character")
	}
	dir, err := StringToDirection(abbrev)
	if err == nil {
		str, err = DirectionToString(dir)
	}
	return
}

//noinspection GoNameStartsWithPackageName
func DirectionToString(d Direction) (str string, err error) {
	str = directionToWord[d]
	if str == "" {
		err = errors.New("Unknown direction")
	}
	return
}

// Given an 'exits' string like "nsed", return a formatted
// string for the player with something like
// "North, South, East, Down"
func ExitsToFormattedString(exits string) (result string, err error) {
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
