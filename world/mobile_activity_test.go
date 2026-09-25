package world

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/watchmud/watchmud/mobile"
	"github.com/watchmud/watchmud/rules"
	"github.com/watchmud/watchmud/spaces"
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
	s.definition = mobile.NewDefinition(

		"id",
		"name",
		"",
		[]string{},
		"desc",
		"room desc",
		25,
		rules.WanderDefinition{
			CanWander:       true,
			CheckFrequency:  time.Minute * 1,
			CheckPercentage: 1.0,
			Style:           rules.WanderFollowPath,
			Path:            []string{"a", "b"},
		},
		10,
		false,
	)
	s.mobileInstance = mobile.NewInstance(s.definition)
}

func (s *MobileActivityTestSuite) Test_getNextDirectionOnPath_Simple() {
	r := spaces.NewTestRoom("a")
	r.Connect(rules.DirectionUp, spaces.NewTestRoom("b"))
	r.DestinationRoom(rules.DirectionUp).Connect(rules.DirectionDown, r)

	// a -> b
	dir, changeDirection, err := getNextDirectionOnPath(s.mobileInstance, r)
	s.Assert().NoError(err)
	s.Assert().Equal(rules.DirectionUp, dir)
	s.Assert().False(changeDirection)

	// b -> a
	dir, changeDirection, err = getNextDirectionOnPath(s.mobileInstance, r.DestinationRoom(rules.DirectionUp))
	s.Assert().NoError(err)
	s.Assert().Equal(rules.DirectionDown, dir)
	s.Assert().True(changeDirection)
}

func (s *MobileActivityTestSuite) Test_getNextDirectionOnPath_FullPath() {
	s.definition.Wandering.Path = []string{"a", "b", "c"}
	// a <-> b <-> c
	a := spaces.NewTestRoom("a")
	b := spaces.NewTestRoom("b")
	c := spaces.NewTestRoom("c")
	a.Connect(rules.DirectionEast, b)
	b.Connect(rules.DirectionWest, a)
	b.Connect(rules.DirectionEast, c)
	c.Connect(rules.DirectionWest, b)

	// a -> b
	dir, changeDirection, err := getNextDirectionOnPath(s.mobileInstance, a)
	s.Assert().NoError(err)
	s.Assert().Equal(rules.DirectionEast, dir)
	s.Assert().False(changeDirection)

	// b -> c
	dir, changeDirection, err = getNextDirectionOnPath(s.mobileInstance, b)
	s.Assert().NoError(err)
	s.Assert().Equal(rules.DirectionEast, dir)
	s.Assert().False(changeDirection)

	// c -> b
	dir, changeDirection, err = getNextDirectionOnPath(s.mobileInstance, c)
	s.Assert().NoError(err)
	s.Assert().Equal(rules.DirectionWest, dir)
	s.Assert().True(changeDirection)

	// b -> a
	// mob needs to be walking back for this to work
	s.mobileInstance.WanderingForward = false
	dir, changeDirection, err = getNextDirectionOnPath(s.mobileInstance, b)
	s.Assert().NoError(err)
	s.Assert().Equal(rules.DirectionWest, dir)
	s.Assert().False(changeDirection) // since we're already walking backwards
}
