package loader

import (
	"fmt"
	"strings"

	"github.com/watchmud/watchmud/mobile"
)

type lootEntry struct {
	// Object is "zone/id", or a bare id for the mob's own zone.
	Object string `json:"object"`
	// Chance is a percent, 1 to 100.
	Chance int `json:"chance"`
}

// mobLoot resolves a mob's loot table against the object definitions. Every
// zone's objects are loaded before any mobs, so a table can name another
// zone's. Anything that doesn't resolve fails startup: a typo in a loot table
// would otherwise be a drop that silently never happens.
func (c *Content) mobLoot(zoneName string, mob mobEntry) ([]mobile.LootEntry, error) {
	var loot []mobile.LootEntry
	for _, l := range mob.Loot {
		zoneId, objectId := zoneName, l.Object
		if z, o, found := strings.Cut(l.Object, "/"); found {
			zoneId, objectId = z, o
		}
		if objectId == "" {
			return nil, fmt.Errorf("mob %s/%s: loot entry names no object", zoneName, mob.Id)
		}
		zone, ok := c.Zones[zoneId]
		if !ok {
			return nil, fmt.Errorf("mob %s/%s: loot %q: zone not found", zoneName, mob.Id, l.Object)
		}
		defn, ok := zone.ObjectDefinitions[objectId]
		if !ok {
			return nil, fmt.Errorf("mob %s/%s: loot %q: object not defined", zoneName, mob.Id, l.Object)
		}
		if l.Chance < 1 || l.Chance > 100 {
			return nil, fmt.Errorf("mob %s/%s: loot %q: chance %d is not 1-100", zoneName, mob.Id, l.Object, l.Chance)
		}
		loot = append(loot, mobile.LootEntry{Object: defn, Chance: l.Chance})
	}
	return loot, nil
}
