package loader

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/watchmud/watchmud/rules"
)

func TestObjectDamage_wielded(t *testing.T) {
	d, err := objectDamage("wrathrock", objectEntry{
		Id:            "sword",
		EquipmentSlot: rules.SlotWield,
		Damage:        "1d8",
	})
	require.NoError(t, err)
	assert.Equal(t, rules.DamageRoll("1d8"), d)
}

func TestObjectDamage_wieldedWithoutDamage(t *testing.T) {
	_, err := objectDamage("wrathrock", objectEntry{
		Id:            "sword",
		EquipmentSlot: rules.SlotWield,
	})
	assert.ErrorContains(t, err, "wrathrock/sword")
}

func TestObjectDamage_onSomethingNotWieldedFails(t *testing.T) {
	_, err := objectDamage("wrathrock", objectEntry{Id: "helm", EquipmentSlot: rules.SlotHead, Damage: "1d4"})
	assert.Error(t, err)
}

func TestObjectDamage_notWieldedIsEmpty(t *testing.T) {
	d, err := objectDamage("wrathrock", objectEntry{Id: "rope"})
	require.NoError(t, err)
	assert.Empty(t, d)
}

func TestObjectDamage_badNotationFails(t *testing.T) {
	_, err := objectDamage("wrathrock", objectEntry{Id: "sword", EquipmentSlot: rules.SlotWield, Damage: "banana 1d6"})
	assert.Error(t, err)
}

func TestMobDamage_defaultsToBareHands(t *testing.T) {
	d, err := mobDamage("wrathrock", mobEntry{Id: "rabbit"})
	require.NoError(t, err)
	assert.Equal(t, rules.BareHands, d)
}

func TestMobDamage_badNotationFails(t *testing.T) {
	_, err := mobDamage("wrathrock", mobEntry{Id: "rabbit", Damage: "lots"})
	assert.ErrorContains(t, err, "wrathrock/rabbit")
}
