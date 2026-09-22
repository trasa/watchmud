package rules

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func testDurabilityTable() DurabilityTable {
	return DurabilityTable{
		Default:        40,
		OnDeathPercent: 10,
		Armor: map[ArmorType]int{
			ArmorTypeCloth: 20,
			ArmorTypePlate: 80,
		},
	}
}

func TestDurabilityTable_armorTypeWins(t *testing.T) {
	table := testDurabilityTable()
	assert.Equal(t, 20, table.MaxFor(ArmorTypeCloth))
	assert.Equal(t, 80, table.MaxFor(ArmorTypePlate))
}

// A type nobody wrote a row for, and gear with no armor type at all, both
// fall back on the one number rather than on nothing.
func TestDurabilityTable_fallsBackToDefault(t *testing.T) {
	table := testDurabilityTable()
	assert.Equal(t, 40, table.MaxFor(ArmorTypeLeather), "no leather row")
	assert.Equal(t, 40, table.MaxFor(ArmorTypeNone), "a knife is not armor")
}

// No durability.json is durability switched off, not every item in the game
// starting out broken.
func TestDurabilityTable_emptyTableIsIndestructible(t *testing.T) {
	var table DurabilityTable
	assert.Equal(t, Indestructible, table.MaxFor(ArmorTypePlate))
	assert.Equal(t, Indestructible, table.MaxFor(ArmorTypeNone))
}

// Dying costs a share of what each piece started at, so the expensive armor
// and the cheap tunic cost the same number of deaths.
func TestDeathLoss_isAShareOfTheMax(t *testing.T) {
	table := testDurabilityTable()
	assert.Equal(t, 8, table.DeathLoss(80))
	assert.Equal(t, 4, table.DeathLoss(40))
	assert.Equal(t, 2, table.DeathLoss(20))
}

// A share that rounds to nothing still costs something, or cheap gear would
// be immortal in precisely the situation the rule is there to punish.
func TestDeathLoss_alwaysCostsAtLeastAPoint(t *testing.T) {
	table := testDurabilityTable()
	assert.Equal(t, 1, table.DeathLoss(5), "10% of 5 rounds to zero")
	assert.Equal(t, 1, table.DeathLoss(1))
}

// Gear that does not wear out does not wear out when you die either.
func TestDeathLoss_indestructibleGearPaysNothing(t *testing.T) {
	assert.Equal(t, 0, testDurabilityTable().DeathLoss(Indestructible))
}

// Content that says nothing about dying makes dying free, the same way
// content that says nothing about durability makes everything indestructible.
func TestDeathLoss_offByDefault(t *testing.T) {
	var table DurabilityTable
	assert.Equal(t, 0, table.DeathLoss(80))

	table = DurabilityTable{Default: 40}
	assert.Equal(t, 0, table.DeathLoss(80))
}
