package world

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// A player named after a mob would make "kill rabbit" ambiguous, and one named
// "all" or "self" would collide with the target grammar.
func TestIsReservedName(t *testing.T) {
	w, err := NewTestWorld()
	require.NoError(t, err)

	for _, name := range []string{
		"All", "Self", "Me", "Corpse", "Someone", // the grammar, and what renders for the unseen
		"Rabbit",          // a mob's name
		"Drone", "Little", // mob aliases
	} {
		assert.True(t, w.IsReservedName(name), name)
	}
	// Walker is a mob in the sample zone, which the test manifest disables:
	// content that isn't loaded reserves nothing
	for _, name := range []string{"Bob", "Rabbits", "Targetdrone", "Walker"} {
		assert.False(t, w.IsReservedName(name), name)
	}
}
