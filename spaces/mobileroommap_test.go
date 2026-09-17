package spaces

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/trasa/watchmud/mobile"
	"github.com/trasa/watchmud/rules"
)

type mobileRoomMapSuite struct {
	suite.Suite
	mob           *mobile.Definition
	instances     []*mobile.Instance
	room          *Room
	mobileRoomMap *MobileRoomMap
}

func TestMobileRoomMapSuite(t *testing.T) {
	suite.Run(t, new(mobileRoomMapSuite))
}

func (s *mobileRoomMapSuite) SetupTest() {
	s.mob = mobile.NewDefinition(
		"id",
		"name",
		"zone",
		[]string{},
		"shortdesc",
		"roomdesc",
		25, rules.WanderDefinition{},
		10,
		false,
	)
	s.instances = append(s.instances, mobile.NewInstance(s.mob))
	s.instances = append(s.instances, mobile.NewInstance(s.mob))

	s.room = NewTestRoom("test")

	s.mobileRoomMap = NewMobileRoomMap()
	s.mobileRoomMap.Add(s.instances[0], s.room)
	s.mobileRoomMap.Add(s.instances[1], s.room)
}

func (s *mobileRoomMapSuite) TestGetAll() {
	result := s.mobileRoomMap.GetAllMobiles()
	s.Assert().Equal(2, len(result))
}
