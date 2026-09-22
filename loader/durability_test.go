package loader

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/trasa/watchmud/rules"
)

func durabilityTable() rules.DurabilityTable {
	return rules.DurabilityTable{
		Default: 40,
		Armor:   map[rules.ArmorType]int{rules.ArmorTypeCloth: 20, rules.ArmorTypePlate: 80},
	}
}

// The table is the point: a builder adding a breastplate writes no number and
// gets what plate is worth.
func TestObjectDurability_fromTheTable(t *testing.T) {
	d, err := objectDurability("wrathrock", objectEntry{Id: "breastplate", ArmorType: rules.ArmorTypePlate}, durabilityTable())
	require.NoError(t, err)
	assert.Equal(t, 80, d)

	d, err = objectDurability("wrathrock", objectEntry{Id: "knife"}, durabilityTable())
	require.NoError(t, err)
	assert.Equal(t, 40, d, "not armor, so the default")
}

func TestObjectDurability_objectOverridesTheTable(t *testing.T) {
	d, err := objectDurability("wrathrock",
		objectEntry{Id: "dwarven_blade", Durability: intp(500)}, durabilityTable())
	require.NoError(t, err)
	assert.Equal(t, 500, d)
}

// Written on purpose, zero is a thing that never wears out -- which is why
// the field is a pointer and not an int.
func TestObjectDurability_explicitZeroIsIndestructible(t *testing.T) {
	d, err := objectDurability("wrathrock",
		objectEntry{Id: "heirloom", ArmorType: rules.ArmorTypePlate, Durability: intp(0)}, durabilityTable())
	require.NoError(t, err)
	assert.Equal(t, rules.Indestructible, d)
}

func TestObjectDurability_negativeIsAnError(t *testing.T) {
	_, err := objectDurability("wrathrock", objectEntry{Id: "oops", Durability: intp(-1)}, durabilityTable())
	assert.ErrorContains(t, err, "negative durability")
}

// the real thing, table and overrides together: the shipped starting kit.
func TestLoadContent_durability(t *testing.T) {
	c, err := LoadContent(os.DirFS("../content"))
	require.NoError(t, err)

	objects := c.Zones["wrathrock"].ObjectDefinitions
	assert.Equal(t, 20, objects["novice_tunic"].MaxDurability, "cloth, from the table")
	assert.Equal(t, 20, objects["travelers_cloak"].MaxDurability, "cloth too")
	assert.Equal(t, 25, objects["training_dagger"].MaxDurability, "its own number")
	assert.Equal(t, 40, objects["waterskin"].MaxDurability, "the default")

	assert.Equal(t, 10, c.Catalog.Durability.OnDeathPercent, "dying costs a tenth of everything worn")
	assert.Equal(t, 2, c.Catalog.Durability.DeathLoss(objects["novice_tunic"].MaxDurability))
}

// Content with no durability.json hands out gear that never wears out, rather
// than gear that is broken on arrival.
func TestLoadContent_noDurabilityFileMeansIndestructible(t *testing.T) {
	c, err := LoadContent(os.DirFS("../testcontent"))
	require.NoError(t, err)

	// testcontent deliberately has no durability.json, which is also what
	// keeps every other suite in the repo playing the game it played before
	// this existed.
	objects := c.Zones["wrathrock"].ObjectDefinitions
	assert.Equal(t, rules.Indestructible, objects["knife"].MaxDurability)
	assert.Equal(t, rules.Indestructible, objects["chain_shirt"].MaxDurability)
}
