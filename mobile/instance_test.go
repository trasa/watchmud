package mobile

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/trasa/watchmud/rules"
)

type instanceSuite struct {
	suite.Suite
	noWalkMob   *Instance
	walkerMob   *Instance
	noChanceMob *Instance
	pathMob     *Instance
}

func TestInstanceSuite(t *testing.T) {
	suite.Run(t, new(instanceSuite))
}

func (s *instanceSuite) SetupTest() {
	s.noWalkMob = NewInstance(NewDefinition("nowalk", "nowalk", "zone", []string{}, "", "",
		25,
		rules.WanderDefinition{
			CanWander: false,
		},
		10,
		false,
	))

	s.walkerMob = NewInstance(NewDefinition("walker", "walker", "zone", []string{}, "", "",
		25,
		rules.WanderDefinition{
			CanWander:       true,
			Style:           rules.WanderRandom,
			CheckFrequency:  time.Minute * 1,
			CheckPercentage: 1.0, // 100 %
		},
		10,
		false,
	))

	s.noChanceMob = NewInstance(NewDefinition("nochance", "nochance", "zone", []string{}, "", "",
		25,
		rules.WanderDefinition{
			CanWander:       true,
			Style:           rules.WanderRandom,
			CheckFrequency:  time.Minute * 1,
			CheckPercentage: 0.0, // <-- 0% chance
		}, 10,
		false,
	))

	s.pathMob = NewInstance(NewDefinition("path", "path", "zone", []string{}, "desc", "room desc",
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
	))
}

func (s *instanceSuite) TestCanWander() {

	s.Assert().False(s.noWalkMob.CanWander())
	s.Assert().False(s.walkerMob.CanWander()) // not time yet

	now := time.Now()
	s.walkerMob.LastWanderingTime = now

	s.Assert().False(s.walkerMob.canWander(now)) // not yet
	s.Assert().False(s.walkerMob.canWander(now.Add(time.Second * 30)))
	s.Assert().False(s.walkerMob.canWander(now.Add(time.Second * 60)))
	s.Assert().True(s.walkerMob.canWander(now.Add(time.Second * 61)))
}

func (s *instanceSuite) TestWanderAlwaysFails() {
	r := rand.New(rand.NewSource(1))

	for i := 0; i < 10; i++ {
		s.Assert().False(s.noChanceMob.checkWanderChance(r)) // always fails
	}
}

func (s *instanceSuite) TestWanderAlwaysSucceeds() {
	r := rand.New(rand.NewSource(1))
	for i := 0; i < 10; i++ {
		s.Assert().True(s.walkerMob.checkWanderChance(r))
	}
}

func (s *instanceSuite) TestPathGetIndex() {
	m := s.pathMob
	idx, err := m.GetIndexOnPath("a")
	s.Assert().NoError(err)
	s.Assert().Equal(0, idx)

	idx, err = m.GetIndexOnPath("b")
	s.Assert().NoError(err)
	s.Assert().Equal(1, idx)

	idx, err = m.GetIndexOnPath("foo")
	s.Assert().Error(err)
	s.Assert().Equal(-1, idx)
}

func (s *instanceSuite) TestInstance_GetIndexOnPath_NoPath() {
	m := s.walkerMob
	idx, err := m.GetIndexOnPath("foo")
	s.Assert().Error(err)
	s.Assert().Equal(-1, idx)
}
