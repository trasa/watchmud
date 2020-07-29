package world

import (
	"errors"
	"strconv"
	"strings"
)

type Target struct {
	Quantity   int    //  "2000 coins" = 2000
	Identifier int    // "2.knife" = 2
	Name       string // "knife" = "knife"
	All        bool
}

// Take the target string and parse it into Target struct
// Example usage:
//      drop <item>
//      drop all.<item>
//      drop all
//      drop <number> coins
//   	drop 2.<item>
func parseTarget(target string) (result Target, err error) {

	parts := strings.Split(target, " ")

	if len(parts) > 2 {
		err = errors.New("TOO_MANY_PARTS_IN_TARGET")
		return
	}
	var item string
	if len(parts) == 2 {
		// "x xyz"
		var quant int
		quant, err = strconv.Atoi(parts[0])
		if err != nil {
			return
		}
		result.Quantity = quant
		item = parts[1]
	} else {
		// "xyz"
		item = parts[0]
	}

	dotparts := strings.Split(item, ".")
	if len(dotparts) == 1 {
		if strings.EqualFold(item, "all") {
			// "all"
			result.All = true
		} else {
			// "foo"
			result.Identifier = 0
			result.Name = item
		}
	} else if len(dotparts) == 2 {
		// "x.foo"
		if strings.EqualFold(dotparts[0], "all") {
			// "all.foo"
			result.All = true
			result.Name = dotparts[1]
		} else {
			// "4.foo"
			var identifier int
			identifier, err = strconv.Atoi(dotparts[0])
			if err != nil {
				return
			}
			result.Identifier = identifier
			result.Name = dotparts[1]
		}
	} else {
		// "x.x.x"
		err = errors.New("TOO_MANY_DOTS")
	}
	return
}
