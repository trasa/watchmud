package server

import (
	"github.com/stretchr/testify/assert"
	"github.com/trasa/watchmud/object"
	"testing"
)

func TestAddInventory_New(t *testing.T) {

	p, _ := NewTestClientPlayer("name")
	p.AddInventory(&object.Instance{
		InstanceId: "id",
		Definition: &object.Definition{
			Id:                  "defnid",
			Categories:          object.CategorySet{object.FOOD: true},
			Name:                "name",
			ShortDescription:    "short desc",
			DescriptionOnGround: "on ground",
		},
	})

	invmap := p.GetInventoryMap()
	assert.Equal(t, 1, len(invmap))
	obj := invmap["id"]
	assert.Equal(t, "id", obj.Id())
	assert.Equal(t, "defnid", obj.(*object.Instance).Definition.Id)
}
