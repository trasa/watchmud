package loader

import (
	"fmt"

	"github.com/watchmud/watchmud/rules"
)

// objectDamage is the dice a wielded object hits with. Anything worn in the wield slot
// has to say, and nothing else may: a weapon with no damage would silently hit like
// bare hands, and damage on a helmet would silently do nothing.
func objectDamage(zoneName string, obj objectEntry) (rules.DamageRoll, error) {
	wielded := obj.EquipmentSlot == rules.SlotWield
	switch {
	case wielded && obj.Damage == "":
		return "", fmt.Errorf("object %s/%s: wielded, but has no damage", zoneName, obj.Id)
	case !wielded && obj.Damage != "":
		return "", fmt.Errorf("object %s/%s: not wielded, but has damage", zoneName, obj.Id)
	case !wielded:
		return "", nil
	}
	d, err := rules.ParseDamageRoll(obj.Damage)
	if err != nil {
		return "", fmt.Errorf("object %s/%s: %w", zoneName, obj.Id, err)
	}
	return d, nil
}

// mobDamage is what the mobfile said, or bare-hands if it said nothing; the same
// "an ordinary creature" defaults a missing ac gets.
func mobDamage(zoneName string, mob mobEntry) (rules.DamageRoll, error) {
	if mob.Damage == "" {
		return rules.BareHands, nil

	}
	d, err := rules.ParseDamageRoll(mob.Damage)
	if err != nil {
		return "", fmt.Errorf("mob %s/%s: %w", zoneName, mob.Id, err)
	}
	return d, nil
}
