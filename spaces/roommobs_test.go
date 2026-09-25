package spaces

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/watchmud/watchmud/mobile"
	"github.com/watchmud/watchmud/rules"
)

type RoomMobsSuite struct {
	suite.Suite
	roomMobs *RoomMobs
	lizard   *mobile.Definition
	rat      *mobile.Definition
}

func TestRoomMobsSuite(t *testing.T) {
	suite.Run(t, new(RoomMobsSuite))
}

func (s *RoomMobsSuite) SetupTest() {
	s.roomMobs = NewRoomMobs()
	s.lizard = newMobDefinition("lizard", []string{"scaly"})
	s.rat = newMobDefinition("rat", nil)
}

func newMobDefinition(name string, aliases []string) *mobile.Definition {
	return mobile.NewDefinition(
		name,
		name,
		"zone",
		aliases,
		"shortdesc",
		"roomdesc",
		25,
		rules.WanderDefinition{},
		10,
		false,
	)
}

// all returns the mobs in the order All() yields them.
func (s *RoomMobsSuite) all() []*mobile.Instance {
	return slices.Collect(s.roomMobs.All())
}

func (s *RoomMobsSuite) TestAddKeepsInsertionOrder() {
	one := mobile.NewInstance(s.lizard)
	two := mobile.NewInstance(s.lizard)
	three := mobile.NewInstance(s.rat)

	s.Require().NoError(s.roomMobs.Add(one))
	s.Require().NoError(s.roomMobs.Add(two))
	s.Require().NoError(s.roomMobs.Add(three))

	s.Equal([]*mobile.Instance{one, two, three}, s.all())
}

// mob1 added, mob2 added, mob1 leaves, mob3 added -> [mob2, mob3]
func (s *RoomMobsSuite) TestRemoveFromMiddleKeepsOrder() {
	one := mobile.NewInstance(s.lizard)
	two := mobile.NewInstance(s.lizard)
	three := mobile.NewInstance(s.rat)

	s.Require().NoError(s.roomMobs.Add(one))
	s.Require().NoError(s.roomMobs.Add(two))
	s.Require().NoError(s.roomMobs.Remove(one))
	s.Require().NoError(s.roomMobs.Add(three))

	s.Equal([]*mobile.Instance{two, three}, s.all())
}

func (s *RoomMobsSuite) TestAddDuplicateIsError() {
	one := mobile.NewInstance(s.lizard)
	s.Require().NoError(s.roomMobs.Add(one))

	s.Error(s.roomMobs.Add(one))
	s.Equal([]*mobile.Instance{one}, s.all())
}

func (s *RoomMobsSuite) TestRemoveMissingIsError() {
	one := mobile.NewInstance(s.lizard)

	s.Error(s.roomMobs.Remove(one))
	s.Empty(s.all())
}

func (s *RoomMobsSuite) TestRemoveEveryMob() {
	one := mobile.NewInstance(s.lizard)
	two := mobile.NewInstance(s.rat)
	s.Require().NoError(s.roomMobs.Add(one))
	s.Require().NoError(s.roomMobs.Add(two))

	s.Require().NoError(s.roomMobs.Remove(two))
	s.Require().NoError(s.roomMobs.Remove(one))

	s.Empty(s.all())
}

// Find returns the earliest-added match, not whichever the map felt like.
func (s *RoomMobsSuite) TestFindReturnsFirstAdded() {
	one := mobile.NewInstance(s.lizard)
	two := mobile.NewInstance(s.lizard)
	s.Require().NoError(s.roomMobs.Add(one))
	s.Require().NoError(s.roomMobs.Add(two))

	found, exists := s.roomMobs.Find("lizard")
	s.Require().True(exists)
	s.Same(one, found)

	// and once the first one leaves, the next in line answers
	s.Require().NoError(s.roomMobs.Remove(one))
	found, exists = s.roomMobs.Find("lizard")
	s.Require().True(exists)
	s.Same(two, found)
}

func (s *RoomMobsSuite) TestFindByAlias() {
	one := mobile.NewInstance(s.lizard)
	s.Require().NoError(s.roomMobs.Add(one))

	found, exists := s.roomMobs.Find("scaly")
	s.Require().True(exists)
	s.Same(one, found)
}

func (s *RoomMobsSuite) TestFindNotFound() {
	s.Require().NoError(s.roomMobs.Add(mobile.NewInstance(s.lizard)))

	found, exists := s.roomMobs.Find("goblin")
	s.False(exists)
	s.Nil(found)
}
