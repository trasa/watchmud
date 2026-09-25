package loader

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/watchmud/watchmud/rules"
)

func intp(i int) *int { return &i }

// A mob file that says nothing about armor class gets the same baseline a
// player with no armor has. It used to get zero, which a d20 can never roll
// under: the mob could not be missed.
func TestMobArmorClass_defaultsToTheBaseline(t *testing.T) {
	ac, err := mobArmorClass("wrathrock", mobEntry{Id: "rabbit"})
	require.NoError(t, err)
	assert.Equal(t, rules.BaseArmorClass, ac)
}

// An explicit number is absolute, not a bonus on top of the baseline: "ac": 10
// is an unarmored creature, and doubling it would make every mob in content
// hittable only on a natural 20.
func TestMobArmorClass_explicitIsAbsolute(t *testing.T) {
	ac, err := mobArmorClass("wrathrock", mobEntry{Id: "drone", AC: intp(10)})
	require.NoError(t, err)
	assert.Equal(t, 10, ac)

	ac, err = mobArmorClass("wrathrock", mobEntry{Id: "knight", AC: intp(18)})
	require.NoError(t, err)
	assert.Equal(t, 18, ac)
}

// Written on purpose, zero is allowed: it means a thing that cannot be
// missed, which is a target dummy, and the pointer is what tells it apart
// from the key being absent.
func TestMobArmorClass_explicitZeroIsKept(t *testing.T) {
	ac, err := mobArmorClass("wrathrock", mobEntry{Id: "dummy", AC: intp(0)})
	require.NoError(t, err)
	assert.Equal(t, 0, ac)
}

func TestMobArmorClass_negativeIsAnError(t *testing.T) {
	_, err := mobArmorClass("wrathrock", mobEntry{Id: "oops", AC: intp(-3)})
	assert.ErrorContains(t, err, "negative ac")
}

// and it holds through a real load
func TestLoadContent_mobArmorClass(t *testing.T) {
	c, err := LoadContent(os.DirFS("../testcontent"))
	require.NoError(t, err)

	mobs := c.Zones["wrathrock"].MobileDefinitions
	assert.Equal(t, 10, mobs["targetDrone"].ArmorClass(), `"ac": 10 is unarmored`)
	assert.Equal(t, rules.BaseArmorClass, mobs["littleDrone"].ArmorClass(), "a mob with no ac gets the baseline")

	// sample's walker declares one too; nothing in content should be sitting
	// at zero by accident.
	for zoneId, zone := range c.Zones {
		for id, defn := range zone.MobileDefinitions {
			assert.Positive(t, defn.ArmorClass(), "mob %s/%s has an armor class of zero: it can never be missed", zoneId, id)
		}
	}
}
