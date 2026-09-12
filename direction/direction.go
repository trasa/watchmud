package direction

import (
	"errors"
	"fmt"
	"strings"
)

// Direction is a compass direction, plus up and down.
//
// The zero value, None, means "unset" and is not a useable direction.
type Direction int

const (
	None Direction = iota
	North
	East
	South
	West
	Up
	Down
)

// ErrUnknownDirection is returned by the parse functions. Callers can test
// for it with errors.Is; the returned error also names the offending input.
var ErrUnknownDirection = errors.New("unknown direction")

// The single source of truth. Indexing by the constants keeps the table and
// the constants from drifting apart, and the compiler sizes the array.
var table = [...]struct {
	name   string
	abbrev string
}{
	None:  {"None!", ""},
	North: {"North", "n"},
	East:  {"East", "e"},
	South: {"South", "s"},
	West:  {"West", "w"},
	Up:    {"Up", "u"},
	Down:  {"Down", "d"},
}

// All is every usable direction, in const order
var All = []Direction{North, East, South, West, Up, Down}

// Valid reports whether d is a usable direction. None is not.
func (d Direction) Valid() bool {
	return d >= North && d <= Down
}

// String returns the display name, e.g. "North"
func (d Direction) String() string {
	if d < 0 || int(d) >= len(table) {
		return fmt.Sprintf("Direction(%d)", int32(d))
	}
	return table[d].name
}

// Abbrev returns the one-letter form, e.g. "n". Empty for None.
func (d Direction) Abbrev() string {
	if !d.Valid() {
		return ""
	}
	return table[d].abbrev
}

// Parse resolves a direction exactly: either the full name or the one-letter
// abbreviation, case-insensitively. Use this for world files, config, and
// anywhere a typo should be an error rather than a guess.
func Parse(s string) (Direction, error) {
	s = strings.TrimSpace(strings.ToLower(s))
	for _, d := range All {
		if s == table[d].abbrev || s == strings.ToLower(table[d].name) {
			return d, nil
		}
	}
	return None, fmt.Errorf("%w: %q", ErrUnknownDirection, s)
}

// ParsePrefix resolves a direction from player input, accepting any unambiguous
// prefix of a direction name: "n", "no", "nort" and "north" all give North.
// Use this for command parsing, never for files.
func ParsePrefix(s string) (Direction, error) {
	s = strings.TrimSpace(strings.ToLower(s))
	if s == "" {
		return None, fmt.Errorf("%w: empty", ErrUnknownDirection)
	}
	for _, d := range All {
		if strings.HasPrefix(strings.ToLower(table[d].name), s) {
			return d, nil
		}
	}
	return None, fmt.Errorf("%w: %q", ErrUnknownDirection, s)
}

// Format renders a set of directions for display:
// "North, South, East, Down", or "None!" when empty.
func Format(dirs []Direction) string {
	if len(dirs) == 0 {
		return table[None].name
	}
	names := make([]string, 0, len(dirs))
	for _, d := range dirs {
		names = append(names, d.String())
	}
	return strings.Join(names, ", ")
}

func (d *Direction) UnmarshalText(text []byte) error {
	v, err := Parse(string(text))
	if err != nil {
		return err
	}
	*d = v
	return nil
}

func (d Direction) MarshalText() ([]byte, error) {
	if !d.Valid() {
		return nil, fmt.Errorf("%w: %d", ErrUnknownDirection, int32(d))
	}
	return []byte(strings.ToLower(table[d].name)), nil
}
