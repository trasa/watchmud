package server

import (
	"github.com/stretchr/testify/assert"
	"github.com/trasa/watchmud/object"
	"github.com/trasa/watchmud/player"
	"testing"
)

func TestAddInventory_New(t *testing.T) {

	p, _ := NewTestClientPlayer("name")
	p.AddInventory(player.InventoryItem{
		Id:               "id",
		Quantity:         3,
		ObjectCategory:   object.FOOD,
		ShortDescription: "short desc",
	})

	invmap := p.GetInventory()
	assert.Equal(t, 1, len(invmap))
	obj := invmap["id"][0]
	assert.Equal(t, 3, obj.Quantity)
}
