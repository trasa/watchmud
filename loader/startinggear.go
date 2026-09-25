package loader

import (
	"fmt"

	"github.com/watchmud/watchmud/rules"
)

// checkStartingGear makes sure every item in the kit is something the world
// can actually hand out, and fails startup when it isn't.
//
// This is a hard failure rather than a skip, unlike the missing definitions
// player.FromRecord tolerates. That one is a save file meeting content that
// moved underneath it, at 2am, for one player; this one is the content
// disagreeing with itself while the server is still starting, and the builder
// who made the typo is right here to read the message.
func (c *Content) checkStartingGear() error {
	// slots claimed by something worn, so two items fighting over the same
	// slot is caught here rather than silently leaving one in inventory.
	claimed := make(map[rules.EquipmentSlot]string, len(c.Catalog.StartingGear))

	for _, item := range c.Catalog.StartingGear {
		zone, ok := c.Zones[item.ZoneId]
		if !ok {
			return fmt.Errorf("starting gear %s/%s: zone not found (is it enabled in the zone manifest?)",
				item.ZoneId, item.DefinitionId)
		}
		defn, ok := zone.ObjectDefinitions[item.DefinitionId]
		if !ok {
			return fmt.Errorf("starting gear %s/%s: object not defined in that zone",
				item.ZoneId, item.DefinitionId)
		}
		if item.Power < 0 {
			return fmt.Errorf("starting gear %s/%s: negative power %d",
				item.ZoneId, item.DefinitionId, item.Power)
		}
		if !item.Equip {
			continue
		}
		if !defn.Wearable() {
			return fmt.Errorf("starting gear %s/%s: marked equip, but it has no equipment_slot",
				item.ZoneId, item.DefinitionId)
		}
		slot := defn.EquipmentSlot
		if prev, dup := claimed[slot]; dup {
			return fmt.Errorf("starting gear %s/%s: slot %s is already worn by %s",
				item.ZoneId, item.DefinitionId, slot, prev)
		}
		claimed[slot] = item.DefinitionId
	}
	return nil
}
