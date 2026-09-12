package spaces

import (
	"testing"
	"uuid"

	"github.com/stretchr/testify/suite"
	"github.com/trasa/watchmud/object"
	"github.com/trasa/watchmud/slot"
)

type RoomInventorySuite struct {
	suite.Suite
	roomInventory *RoomInventory
	defn          *object.Definition
	inst          *object.Instance
	instTwo       *object.Instance
}

func TestRoomInventorySuite(t *testing.T) {
	suite.Run(t, new(RoomInventorySuite))
}

func (suite *RoomInventorySuite) SetupTest() {
	suite.roomInventory = NewRoomInventory()
	suite.defn = object.NewDefinition("id", "name", "zoneid", object.Other, []string{}, "short desc", "on ground", slot.None)
	suite.inst = object.NewInstance(uuid.New(), suite.defn)
	suite.instTwo = object.NewInstance(uuid.New(), suite.defn)

	suite.Require().NoError(suite.roomInventory.Add(suite.inst))
	suite.Require().NoError(suite.roomInventory.Add(suite.instTwo))
}

func (suite *RoomInventorySuite) TestRoomInventory_AddMany() {

	instances := suite.roomInventory.Name("name")
	suite.Assert().NotEmpty(instances)

	all := suite.roomInventory.GetAll()
	suite.Assert().Equal(2, len(all))

	retone, exists := suite.roomInventory.InstanceId(suite.inst.Id)
	suite.Assert().True(exists)
	suite.Assert().Equal(retone, suite.inst)

	rettwo, exists := suite.roomInventory.InstanceId(suite.instTwo.Id)
	suite.Assert().True(exists)
	suite.Assert().Equal(rettwo, suite.instTwo)

	nothing, exists := suite.roomInventory.InstanceId(uuid.New())
	suite.Assert().False(exists)
	suite.Assert().Nil(nothing)
}

func (suite *RoomInventorySuite) TestRoomInventory_Remove() {

	suite.Assert().NoError(suite.roomInventory.Remove(suite.inst))

	suite.Assert().Equal(1, len(suite.roomInventory.GetAll()))

	ret, exists := suite.roomInventory.InstanceId(suite.inst.Id)
	suite.Assert().False(exists)
	suite.Assert().Nil(ret)
}
