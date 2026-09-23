package player

import (
	"testing"
	"uuid"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/trasa/watchmud/object"
	"github.com/trasa/watchmud/rules"
)

func TestWeaponDamageRoll(t *testing.T) {
	p := NewTestPlayer(uuid.New(), "dood", nil)
	assert.Equal(t, string(rules.BareHands), p.WeaponDamageRoll(), "nothing wielded")

	d := object.NewDefinition(
		"sword",
		"sword",
		"wrathrock",
		rules.ObjectCategoryWeapon,
		nil,
		"a sword",
		"A sword is here.",
		rules.SlotWield,
		rules.ArmorTypeNone)
	d.Damage = "1d8"
	d.MaxDurability = 1
	sword := object.NewInstance(uuid.New(), d)
	p.Equipment().Equip(rules.SlotWield, sword)
	assert.Equal(t, "1d8", p.WeaponDamageRoll(), "sword wielded")

	sword.Damage(1)
	require.True(t, sword.Broken())
	assert.Equal(t, string(rules.BareHands), p.WeaponDamageRoll(), "broken sword")
}
