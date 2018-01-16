package server

import (
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/trasa/watchmud/object"
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
		object.FOOD, []string{}, "short desc", "in room")
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

func (suite *ClientPlayerSuite) TestSetEquippedPrimaryWeapon() {
	weaponInst := &object.Instance{
		InstanceId: uuid.NewV4(),
		Definition: object.NewDefinition("weapon", "weapon", "zone", object.WEAPON, []string{}, "weapon", "weapon"),
	}
	suite.player.AddInventory(weaponInst)

	// success
	assert.NoError(suite.T(), suite.player.SetEquippedPrimaryWeapon(weaponInst))
}

func (suite *ClientPlayerSuite) TestCantEquipYouDontHaveOne() {
	youdonthaveoneInst := &object.Instance{
		InstanceId: uuid.NewV4(),
		Definition: object.NewDefinition("nothere", "nothere", "zone", object.WEAPON, []string{}, "youdonthaveone", "youdonthaveone"),
	}
	assert.IsType(suite.T(), &object.InstanceNotFoundError{}, suite.player.SetEquippedPrimaryWeapon(youdonthaveoneInst))
}

func (suite *ClientPlayerSuite) TestNotEquipableWeapon() {
	cantequipthat := &object.Instance{
		InstanceId: uuid.NewV4(),
		Definition: object.NewDefinition("treasure", "treasure", "zone", object.TREASURE, []string{}, "treasure", "treasure"),
	}
	suite.player.AddInventory(cantequipthat)

	// that isn't a weapon
	assert.IsType(suite.T(), &object.InstanceNotWeaponError{}, suite.player.SetEquippedPrimaryWeapon(cantequipthat))
}
