package player

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type LocationTestSuite struct {
	suite.Suite
}

func TestLocationTestSuite(t *testing.T) {
	suite.Run(t, new(LocationTestSuite))
}

func (s *LocationTestSuite) Test_Null() {
	l := NewLocation(nil, nil)

	s.Assert().Empty(l.RoomId)
	s.Assert().Empty(l.ZoneId)
}

func (s *LocationTestSuite) Test_NotNull() {
	z := "zoneid"
	r := "roomId"
	l := NewLocation(&z, &r)

	s.Assert().Equal(z, l.ZoneId)
	s.Assert().Equal(r, l.RoomId)
}
