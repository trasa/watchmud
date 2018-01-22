package server

import (
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/trasa/watchmud/object"
	"github.com/trasa/watchmud/slot"
	"testing"
)

type ClientPlayerSuite struct {
	suite.Suite
	player *ClientPlayer
}

func TestClientPlayerSuite(t *testing.T) {
	suite.Run(t, new(ClientPlayerSuite))
}

func (suite *ClientPlayerSuite) SetupTest() {
	suite.player, _ = NewTestClientPlayer("name")
}

func (suite *ClientPlayerSuite) TestAddInventory_New() {

	defnPtr := object.NewDefinition("defnid", "name", "zone",
		object.FOOD, []string{}, "short desc", "in room", slot.None)
	instPtr := &object.Instance{
		InstanceId: uuid.NewV4(),
		Definition: defnPtr,
	}

	suite.player.AddInventory(instPtr)

	invs := suite.player.GetAllInventory()
	assert.Equal(suite.T(), 1, len(invs))
	obj := invs[0]
	assert.Equal(suite.T(), instPtr.Id(), obj.Id())
	assert.Equal(suite.T(), "defnid", obj.Definition.Identifier())
}
