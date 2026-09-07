package player

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type LocationTestSuite struct {
	suite.Suite
}

func TestLocationTestSuite(t *testing.T) {
	suite.Run(t, new(LocationTestSuite))
}

func (s *LocationTestSuite) TestNewLocation() {
	z := "zoneid"
	r := "roomId"
	l := NewLocation(z, r)

	s.Assert().Equal(z, l.ZoneId)
	s.Assert().Equal(r, l.RoomId)
}
