package spaces

import (
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/suite"
	"github.com/trasa/watchmud/object"
	"github.com/trasa/watchmud/slot"
	"testing"
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
	suite.inst = object.NewInstance(suite.defn)
	suite.instTwo = object.NewInstance(suite.defn)
	suite.roomInventory.Add(suite.inst)
	suite.roomInventory.Add(suite.instTwo)
}

func (suite *RoomInventorySuite) TestRoomInventory_AddMany() {

	inst, exists := suite.roomInventory.GetByName("name")
	suite.Assert().True(exists)
	suite.Assert().NotNil(inst)

	all := suite.roomInventory.GetAll()
	suite.Assert().Equal(2, len(all))

	retone, exists := suite.roomInventory.GetByInstanceId(suite.inst.InstanceId)
	suite.Assert().True(exists)
	suite.Assert().Equal(retone, suite.inst)

	rettwo, exists := suite.roomInventory.GetByInstanceId(suite.instTwo.InstanceId)
	suite.Assert().True(exists)
	suite.Assert().Equal(rettwo, suite.instTwo)

	nothing, exists := suite.roomInventory.GetByInstanceId(uuid.Must(uuid.NewV4()))
	suite.Assert().False(exists)
	suite.Assert().Nil(nothing)
}

func (suite *RoomInventorySuite) TestRoomInventory_Remove() {

	suite.Assert().NoError(suite.roomInventory.Remove(suite.inst))

	suite.Assert().Equal(1, len(suite.roomInventory.GetAll()))

	ret, exists := suite.roomInventory.GetByInstanceId(suite.inst.InstanceId)
	suite.Assert().False(exists)
	suite.Assert().Nil(ret)
}
