package slot

import (
	"fmt"
	"strings"
)

type Location int32

//go:generate stringer -type=Location
const (
	None  Location = iota
	Wield          // item can be wielded (i.e. weapon)
	Hold           // item can be held ('hold' command)
	Head
	Neck
	Body
	About_Body // like a cloak
	Legs
	Feet
	Arms
	Wrist
	Hands
	Fingers
	Waist
)

func StringToLocation(name string) (Location, error) {
	if name == "" {
		return None, nil
	}
	name = strings.ToUpper(name)
	stridx := strings.Index(strings.ToUpper(_Location_name), name)
	if stridx < 0 {
		return None, fmt.Errorf("'%s' not found", name)
	}

	for pos, catidx := range _Location_index {
		if stridx == int(catidx) {
			return Location(pos), nil
		}
	}
	// shouldn't happen?
	return None, fmt.Errorf("could not find index %d for '%s'", stridx, name)
}
