package loader

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/watchmud/watchmud/rules"
	"github.com/watchmud/watchmud/spaces"
)

// walk follows exits from a room, failing the test on the first one that
// isn't there -- a broken link fails the build instead of a player.
func walk(t *testing.T, c *Content, zoneId, roomId string, dirs ...rules.Direction) *spaces.Room {
	t.Helper()
	room, err := c.Room(zoneId, roomId)
	require.NoError(t, err)
	for _, d := range dirs {
		next := room.DestinationRoom(d)
		require.NotNil(t, next, "no exit %s from %s/%s", d, room.Zone.Id, room.Id)
		room = next
	}
	return room
}

// The newbie country south of Wrathrock and the barrow beneath it, and the
// boss who is meant to be out of reach of anyone on their own. See LEVELS.md.
func TestLoadContent_hollowfieldAndBarrow(t *testing.T) {
	c, err := LoadContent(os.DirFS("../content"))
	require.NoError(t, err)

	hollowfield, barrow := c.Zones["hollowfield"], c.Zones["barrow"]
	require.NotNil(t, hollowfield)
	require.NotNil(t, barrow)
	assert.Equal(t, rules.PowerBand{Min: 1, Max: 5}, hollowfield.Power)
	assert.Equal(t, rules.PowerBand{Min: 6, Max: 10}, barrow.Power)

	// newbie mobs sit inside the band
	for id, m := range hollowfield.MobileDefinitions {
		assert.GreaterOrEqual(t, m.Power, hollowfield.Power.Min, id)
		assert.LessOrEqual(t, m.Power, hollowfield.Power.Max, id)
	}
	// and the boss sits well above his own
	king := barrow.MobileDefinitions["barrow_king"]
	require.NotNil(t, king)
	assert.Greater(t, king.Power, barrow.Power.Max)
	assert.True(t, king.HasFlag("aggressive"))
}

// From Wrathrock's southern path all the way down to the throne, and back up.
func TestLoadContent_southToTheThroneRoom(t *testing.T) {
	c, err := LoadContent(os.DirFS("../content"))
	require.NoError(t, err)

	N, S, E, W := rules.DirectionNorth, rules.DirectionSouth, rules.DirectionEast, rules.DirectionWest
	U, D := rules.DirectionUp, rules.DirectionDown

	oak := walk(t, c, "wrathrock", "south_trail", S, S, S, S, S, E)
	assert.Equal(t, "overturned_oak", oak.Id)

	throne := walk(t, c, "hollowfield", "overturned_oak", D, D, S, E, S, D)
	assert.Equal(t, "throne_room", throne.Id)

	back := walk(t, c, "barrow", "throne_room", U, N, W, N, U, U, W, N, N, N, N, N)
	assert.Equal(t, "wrathrock", back.Zone.Id)
	assert.Equal(t, "south_trail", back.Id)
}

// Every exit in the two new zones has a way back the other way, so nobody
// walks into a room they can't walk out of.
func TestLoadContent_hollowfieldExitsAreTwoWay(t *testing.T) {
	c, err := LoadContent(os.DirFS("../content"))
	require.NoError(t, err)

	opposite := map[rules.Direction]rules.Direction{
		rules.DirectionNorth: rules.DirectionSouth, rules.DirectionSouth: rules.DirectionNorth,
		rules.DirectionEast: rules.DirectionWest, rules.DirectionWest: rules.DirectionEast,
		rules.DirectionUp: rules.DirectionDown, rules.DirectionDown: rules.DirectionUp,
	}
	for _, zoneId := range []string{"hollowfield", "barrow"} {
		for id, room := range c.Zones[zoneId].Rooms {
			for _, d := range rules.AllUsableDirections {
				dest := room.DestinationRoom(d)
				if dest == nil {
					continue
				}
				assert.Same(t, room, dest.DestinationRoom(opposite[d]),
					"%s/%s goes %s to %s/%s, which doesn't lead back", zoneId, id, d, dest.Zone.Id, dest.Id)
			}
		}
	}
}
