package loader

import (
	"fmt"

	"github.com/watchmud/watchmud/rules"
)

// zonePowerBand is what the manifest said, or the bottom band if it said
// nothing.
func zonePowerBand(m zoneManifestEntry) (rules.PowerBand, error) {
	if m.Power == nil {
		return rules.PowerBand{}, nil
	}
	b := *m.Power
	if b.Min < 0 || b.Max < b.Min {
		return rules.PowerBand{}, fmt.Errorf("zone %s: bad power band %d-%d", m.Id, b.Min, b.Max)
	}
	return b, nil
}

// mobPower is what the mob file said, or the bottom of its zone's band if it
// said nothing -- an ordinary inhabitant, the same way a missing ac is an
// ordinary target. An explicit number may sit outside the band; a boss does.
func mobPower(zoneName string, mob mobEntry, band rules.PowerBand) (int, error) {
	if mob.Power == nil {
		return band.Min, nil
	}
	if *mob.Power < 0 {
		return 0, fmt.Errorf("mob %s/%s: negative power %d", zoneName, mob.Id, *mob.Power)
	}
	return *mob.Power, nil
}
