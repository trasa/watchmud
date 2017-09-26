package world

import (
	"github.com/stretchr/testify/assert"
	"github.com/trasa/watchmud/player"
	"testing"
)

func TestPlayerRoomMap_GetPlayers(t *testing.T) {
	m := NewPlayerRoomMap()
	bob := player.NewTestPlayer("bob")
	alice := player.NewTestPlayer("alice")

	northRoom := NewTestRoom("north")

	m.Add(bob, northRoom)
	m.Add(alice, northRoom)

	assert.Equal(t, 2, len(m.GetPlayers(northRoom)))

	m.Remove(bob)

	assert.Equal(t, 1, len(m.GetPlayers(northRoom)))

}
