package world

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPlayerRoomMap_GetPlayers(t *testing.T) {
	m := NewPlayerRoomMap()
	bob := NewTestPlayer("bob")
	alice := NewTestPlayer("alice")

	northRoom := NewTestRoom("north")

	m.Add(bob, &northRoom)
	m.Add(alice, &northRoom)

	assert.Equal(t, 2, len(m.GetPlayers(&northRoom)))

	m.Remove(bob)

	assert.Equal(t, 1, len(m.GetPlayers(&northRoom)))

}
