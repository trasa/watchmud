package loader

import "github.com/trasa/watchmud/spaces"

type WorldBuilder struct {
	zones map[string]*spaces.Zone
}

func BuildWorld() map[string]*spaces.Zone {
	worldBuilder := WorldBuilder{
		make(map[string]*spaces.Zone),
	}

	worldBuilder.loadZoneManifest()

	return worldBuilder.zones
}

// Retrieve the zone manifest; prepare the zone objects to be
// populated by rooms, objects, mobiles (but don't process the
// zone commands yet)
func (worldBuilder *WorldBuilder) loadZoneManifest() {

	// here, we'd look up something from the database, or something.
	sampleZone := spaces.NewZone("void", "void zone")
	worldBuilder.zones[sampleZone.Id] = sampleZone
}
