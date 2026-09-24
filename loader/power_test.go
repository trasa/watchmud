package loader

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/trasa/watchmud/rules"
)

// A zone that says nothing about power is the bottom band: everything in it
// is power 0, which is what everything was before power existed.
func TestZonePowerBand_defaultsToZero(t *testing.T) {
	band, err := zonePowerBand(zoneManifestEntry{Id: "void"})
	require.NoError(t, err)
	assert.Equal(t, rules.PowerBand{}, band)
}

func TestZonePowerBand_explicit(t *testing.T) {
	band, err := zonePowerBand(zoneManifestEntry{Id: "caves", Power: &rules.PowerBand{Min: 10, Max: 15}})
	require.NoError(t, err)
	assert.Equal(t, rules.PowerBand{Min: 10, Max: 15}, band)
}

func TestZonePowerBand_invalid(t *testing.T) {
	_, err := zonePowerBand(zoneManifestEntry{Id: "caves", Power: &rules.PowerBand{Min: 15, Max: 10}})
	assert.ErrorContains(t, err, "caves")

	_, err = zonePowerBand(zoneManifestEntry{Id: "caves", Power: &rules.PowerBand{Min: -1, Max: 10}})
	assert.ErrorContains(t, err, "caves")
}

// A mob that doesn't say is an ordinary inhabitant of its zone: the bottom of
// the band, so a builder who forgets makes the easy version, not the boss.
func TestMobPower_defaultsToTheBottomOfTheZoneBand(t *testing.T) {
	p, err := mobPower("caves", mobEntry{Id: "bat"}, rules.PowerBand{Min: 10, Max: 15})
	require.NoError(t, err)
	assert.Equal(t, 10, p)
}

// An explicit number is kept, even outside the band -- that's what a boss is.
func TestMobPower_explicitIsKept(t *testing.T) {
	p, err := mobPower("caves", mobEntry{Id: "bat", Power: intp(12)}, rules.PowerBand{Min: 10, Max: 15})
	require.NoError(t, err)
	assert.Equal(t, 12, p)

	p, err = mobPower("caves", mobEntry{Id: "dragon", Power: intp(25)}, rules.PowerBand{Min: 10, Max: 15})
	require.NoError(t, err)
	assert.Equal(t, 25, p)
}

func TestMobPower_negativeIsAnError(t *testing.T) {
	_, err := mobPower("caves", mobEntry{Id: "oops", Power: intp(-1)}, rules.PowerBand{})
	assert.ErrorContains(t, err, "negative power")
}

// and through a real load
func TestLoadContent_power(t *testing.T) {
	c, err := LoadContent(os.DirFS("../testcontent"))
	require.NoError(t, err)

	wrathrock := c.Zones["wrathrock"]
	assert.Equal(t, rules.PowerBand{Min: 1, Max: 5}, wrathrock.Power)
	assert.Equal(t, 3, wrathrock.MobileDefinitions["targetDrone"].Power, "declared")
	assert.Equal(t, 1, wrathrock.MobileDefinitions["littleDrone"].Power, "bottom of the band")

	assert.Equal(t, rules.PowerBand{}, c.Zones["void"].Power, "no band in the manifest")
}
