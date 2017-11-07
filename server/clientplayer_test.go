package server

import (
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/trasa/watchmud/object"
	"testing"
)

func TestAddInventory_New(t *testing.T) {

	p, _ := NewTestClientPlayer("name")
	defnPtr := object.NewDefinition("defnid", "name", "zone",
		object.FOOD, []string{}, "short desc", "in room")
	instPtr := &object.Instance{
		InstanceId: uuid.NewV4(),
		Definition: defnPtr,
	}

	p.AddInventory(instPtr)

	invs := p.GetAllInventory()
	assert.Equal(t, 1, len(invs))
	obj := invs[0]
	assert.Equal(t, instPtr.Id(), obj.Id())
	assert.Equal(t, "defnid", obj.Definition.Identifier())
}
