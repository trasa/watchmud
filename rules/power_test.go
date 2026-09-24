package rules

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPowerDelta_isClamped(t *testing.T) {
	assert.Equal(t, 0, PowerDelta(5, 5))
	assert.Equal(t, 3, PowerDelta(8, 5))
	assert.Equal(t, -3, PowerDelta(5, 8))
	assert.Equal(t, PowerDeltaClamp, PowerDelta(40, 1))
	assert.Equal(t, -PowerDeltaClamp, PowerDelta(1, 40))
}

// Half a point per point of difference, symmetric either side of zero so
// being three above is worth exactly what being three below costs.
func TestPowerHitModifier(t *testing.T) {
	assert.Equal(t, 0, PowerHitModifier(5, 5))
	assert.Equal(t, 0, PowerHitModifier(6, 5))
	assert.Equal(t, 1, PowerHitModifier(7, 5))
	assert.Equal(t, -1, PowerHitModifier(5, 7))
	assert.Equal(t, -1, PowerHitModifier(5, 8))
	assert.Equal(t, 5, PowerHitModifier(99, 0))
	assert.Equal(t, -5, PowerHitModifier(0, 99))
}

// Five percent per point, rounded to nearest.
func TestPowerDamage(t *testing.T) {
	assert.Equal(t, 10, PowerDamage(10, 5, 5))
	assert.Equal(t, 15, PowerDamage(10, 15, 5), "+10: half again")
	assert.Equal(t, 5, PowerDamage(10, 5, 15), "-10: half")
	assert.Equal(t, 15, PowerDamage(10, 99, 0), "clamped")
	assert.Equal(t, 11, PowerDamage(10, 7, 5), "10 * 1.10")
	assert.Equal(t, 9, PowerDamage(9, 6, 5), "9 * 1.05 = 9.45")
	assert.Equal(t, 10, PowerDamage(9, 7, 5), "9 * 1.10 = 9.9")
}

// A blow that landed does at least a point, however far outclassed.
func TestPowerDamage_aHitIsAtLeastOne(t *testing.T) {
	assert.Equal(t, 1, PowerDamage(1, 0, 15))
	assert.Equal(t, 0, PowerDamage(0, 15, 0), "but no damage stays no damage")
}

// The loot bump, on a d100 roll of 0-99: 2% +2, the next 10% +1.
func TestLootPowerBump(t *testing.T) {
	assert.Equal(t, 2, LootPowerBump(0))
	assert.Equal(t, 2, LootPowerBump(1))
	assert.Equal(t, 1, LootPowerBump(2))
	assert.Equal(t, 1, LootPowerBump(11))
	assert.Equal(t, 0, LootPowerBump(12))
	assert.Equal(t, 0, LootPowerBump(99))
}
