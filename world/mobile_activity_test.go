package world

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/trasa/watchmud-message/direction"
	"github.com/trasa/watchmud/mobile"
	"github.com/trasa/watchmud/spaces"
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
	r.Set(direction.Up, spaces.NewTestRoom("b"))
	r.Get(direction.Up).Set(direction.Down, r)

	// a -> b
	dir, changeDirection, err := getNextDirectionOnPath(s.mobileInstance, r)
	s.Assert().NoError(err)
	s.Assert().Equal(direction.Up, dir)
	s.Assert().False(changeDirection)

	// b -> a
	dir, changeDirection, err = getNextDirectionOnPath(s.mobileInstance, r.Get(direction.Up))
	s.Assert().NoError(err)
	s.Assert().Equal(direction.Down, dir)
	s.Assert().True(changeDirection)
}

func (s *MobileActivityTestSuite) Test_getNextDirectionOnPath_FullPath() {
	s.definition.Wandering.Path = []string{"a", "b", "c"}
	// a <-> b <-> c
	a := spaces.NewTestRoom("a")
	b := spaces.NewTestRoom("b")
	c := spaces.NewTestRoom("c")
	a.Set(direction.East, b)
	b.Set(direction.West, a)
	b.Set(direction.East, c)
	c.Set(direction.West, b)

	// a -> b
	dir, changeDirection, err := getNextDirectionOnPath(s.mobileInstance, a)
	s.Assert().NoError(err)
	s.Assert().Equal(direction.East, dir)
	s.Assert().False(changeDirection)

	// b -> c
	dir, changeDirection, err = getNextDirectionOnPath(s.mobileInstance, b)
	s.Assert().NoError(err)
	s.Assert().Equal(direction.East, dir)
	s.Assert().False(changeDirection)

	// c -> b
	dir, changeDirection, err = getNextDirectionOnPath(s.mobileInstance, c)
	s.Assert().NoError(err)
	s.Assert().Equal(direction.West, dir)
	s.Assert().True(changeDirection)

	// b -> a
	// mob needs to be walking back for this to work
	s.mobileInstance.WanderingForward = false
	dir, changeDirection, err = getNextDirectionOnPath(s.mobileInstance, b)
	s.Assert().NoError(err)
	s.Assert().Equal(direction.West, dir)
	s.Assert().False(changeDirection) // since we're already walking backwards
}
