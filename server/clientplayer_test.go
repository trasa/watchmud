package server

import (
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/suite"
	"github.com/trasa/watchmud-message/slot"
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
		object.Food, []string{}, "short desc", "in room", slot.None)
	instPtr := &object.Instance{
		InstanceId: uuid.NewV4(),
		Definition: defnPtr,
	}

	suite.player.AddInventory(instPtr)

	invs := suite.player.GetAllInventory()
	suite.Assert().Equal(1, len(invs))
	obj := invs[0]
	suite.Assert().Equal(instPtr.Id(), obj.Id())
	suite.Assert().Equal("defnid", obj.Definition.Identifier())
}

func (suite *ClientPlayerSuite) TestMeleeDamage() {
	startingHealth := suite.player.GetCurrentHealth()

	isDead := suite.player.TakeMeleeDamage(5)

	suite.Assert().False(isDead)
	suite.Assert().Equal(startingHealth-5, suite.player.GetCurrentHealth())
}

func (suite *ClientPlayerSuite) TestFatalMeleeDamage() {
	startingHealth := suite.player.GetCurrentHealth()

	isDead := suite.player.TakeMeleeDamage(startingHealth)

	suite.Assert().True(isDead)
	suite.Assert().Equal(0, suite.player.GetCurrentHealth())
}

func (suite *ClientPlayerSuite) TestOverwhelminglyFatalMeleeDamage() {
	startingHealth := suite.player.GetCurrentHealth()

	isDead := suite.player.TakeMeleeDamage(startingHealth * 2)

	suite.Assert().True(isDead)
	suite.Assert().True(suite.player.curHealth < 0)
}

func (suite *ClientPlayerSuite) TestIsDead() {
	suite.Assert().False(suite.player.IsDead())
	suite.player.curHealth = 0
	suite.Assert().True(suite.player.IsDead())
}
