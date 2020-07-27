package world

import (
	"github.com/stretchr/testify/suite"
	"github.com/trasa/watchmud-message/direction"
	"github.com/trasa/watchmud/mobile"
	"github.com/trasa/watchmud/spaces"
	"testing"
	"time"
)

type MobileActivityTestSuite struct {
	suite.Suite
	definition     *mobile.Definition
	mobileInstance *mobile.Instance
}

func TestMobileActivityTestSuite(t *testing.T) {
	suite.Run(t, new(MobileActivityTestSuite))
}

func (s *MobileActivityTestSuite) SetupTest() {
	s.definition = mobile.NewDefinition("id", "name", "", []string{}, "desc", "room desc", 25, mobile.WanderingDefinition{
		CanWander:       true,
		CheckFrequency:  time.Minute * 1,
		CheckPercentage: 1.0,
		Style:           mobile.WANDER_FOLLOW_PATH,
		Path:            []string{"a", "b"},
	}, 10)
	s.mobileInstance = mobile.NewInstance(s.definition)
}

func (s *MobileActivityTestSuite) Test_getNextDirectionOnPath_Simple() {
	r := spaces.NewTestRoom("a")
	r.Set(direction.UP, spaces.NewTestRoom("b"))
	r.Get(direction.UP).Set(direction.DOWN, r)

	// a -> b
	dir, changeDirection, err := getNextDirectionOnPath(s.mobileInstance, r)
	s.Assert().NoError(err)
	s.Assert().Equal(direction.UP, dir)
	s.Assert().False(changeDirection)

	// b -> a
	dir, changeDirection, err = getNextDirectionOnPath(s.mobileInstance, r.Get(direction.UP))
	s.Assert().NoError(err)
	s.Assert().Equal(direction.DOWN, dir)
	s.Assert().True(changeDirection)
}

func (s *MobileActivityTestSuite) Test_getNextDirectionOnPath_FullPath() {
	s.definition.Wandering.Path = []string{"a", "b", "c"}
	// a <-> b <-> c
	a := spaces.NewTestRoom("a")
	b := spaces.NewTestRoom("b")
	c := spaces.NewTestRoom("c")
	a.Set(direction.EAST, b)
	b.Set(direction.WEST, a)
	b.Set(direction.EAST, c)
	c.Set(direction.WEST, b)

	// a -> b
	dir, changeDirection, err := getNextDirectionOnPath(s.mobileInstance, a)
	s.Assert().NoError(err)
	s.Assert().Equal(direction.EAST, dir)
	s.Assert().False(changeDirection)

	// b -> c
	dir, changeDirection, err = getNextDirectionOnPath(s.mobileInstance, b)
	s.Assert().NoError(err)
	s.Assert().Equal(direction.EAST, dir)
	s.Assert().False(changeDirection)

	// c -> b
	dir, changeDirection, err = getNextDirectionOnPath(s.mobileInstance, c)
	s.Assert().NoError(err)
	s.Assert().Equal(direction.WEST, dir)
	s.Assert().True(changeDirection)

	// b -> a
	// mob needs to be walking back for this to work
	s.mobileInstance.WanderingForward = false
	dir, changeDirection, err = getNextDirectionOnPath(s.mobileInstance, b)
	s.Assert().NoError(err)
	s.Assert().Equal(direction.WEST, dir)
	s.Assert().False(changeDirection) // since we're already walking backwards
}
