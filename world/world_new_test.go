package world

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/trasa/watchmud/loader"
	"github.com/trasa/watchmud/rules"
	"github.com/trasa/watchmud/spaces"
	"github.com/trasa/watchmud/zonereset"
)

func TestWorldNew(t *testing.T) {
	settings := loader.Settings{
		StartZone: "z",
		StartRoom: "r",
		VoidZone:  "z",
		VoidRoom:  "r",
	}
	var species []*rules.Species
	var classes []*rules.Class
	catalog, err := rules.NewCatalog(species, classes)
	assert.NoError(t, err)
	content := loader.NewContent(&settings, catalog)
	content.Zones["z"] = spaces.NewZone("z", "Z", zonereset.NEVER, time.Duration(0))
	content.Zones["z"].AddRoom(spaces.NewRoom(content.Zones["z"], "r", "R", "Description"))

	w, err := New(content)
	assert.NoError(t, err)
	assert.NotNil(t, w)
}
