package rules

import (
	"fmt"
	"regexp"
)

// DamageRoll is dice notation for how hard something hits: "1d6", "2d4+1".
// Content only makes one through ParseDamageRoll, so a DamageRoll that
// came through the loader is known to roll.
type DamageRoll string

// BareHands is what you hit with when you're wielding nothing, or
// nothing that still works, and what a mob hits with when its file
// doesn't say. A placeholder, see LEVELS.md "Tuning Placeholders"
const BareHands DamageRoll = "1d2"

var damageRollPattern = regexp.MustCompile(`^[1-9][0-9]*d[1-9][0-9]*([+-][0-9]+)?$`)

// ParseDamageRoll accepts NdM with an optional +K or -K, and nothing else.
// Stricter than the dice library on purpose: that one finds a roll anywhere
// in the string, so "banana 1d6" would roll 1d6 and nobody would notice.
func ParseDamageRoll(s string) (DamageRoll, error) {
	if !damageRollPattern.MatchString(s) {
		return "", fmt.Errorf("damage %q: want dice like 1d6 or 2d4+1", s)
	}
	return DamageRoll(s), nil
}
