package rules

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type AbilitiesTestSuite struct {
	suite.Suite
	abilities Abilities
}

func TestAbilitiesTestSuite(t *testing.T) {
	suite.Run(t, new(AbilitiesTestSuite))
}

func (s *AbilitiesTestSuite) SetupTest() {
	s.abilities = Abilities{}
}

func (s *AbilitiesTestSuite) Test_FillEmptyScoreByPriority() {
	s.abilities = s.abilities.fillEmptyScoreByPriority(18)
	s.Assert().Equal(18, s.abilities.Con)
	s.Assert().Equal(0, s.abilities.Dex)
	s.Assert().Equal(0, s.abilities.Str)
	s.Assert().Equal(0, s.abilities.Int)
	s.Assert().Equal(0, s.abilities.Wis)
	s.Assert().Equal(0, s.abilities.Cha)

	s.abilities = s.abilities.fillEmptyScoreByPriority(17)
	s.Assert().Equal(18, s.abilities.Con)
	s.Assert().Equal(17, s.abilities.Dex)
	s.Assert().Equal(0, s.abilities.Str)
	s.Assert().Equal(0, s.abilities.Int)
	s.Assert().Equal(0, s.abilities.Wis)
	s.Assert().Equal(0, s.abilities.Cha)

	s.abilities = s.abilities.fillEmptyScoreByPriority(16)
	s.Assert().Equal(18, s.abilities.Con)
	s.Assert().Equal(17, s.abilities.Dex)
	s.Assert().Equal(16, s.abilities.Str)
	s.Assert().Equal(0, s.abilities.Int)
	s.Assert().Equal(0, s.abilities.Wis)
	s.Assert().Equal(0, s.abilities.Cha)

	s.abilities = s.abilities.fillEmptyScoreByPriority(15)
	s.Assert().Equal(18, s.abilities.Con)
	s.Assert().Equal(17, s.abilities.Dex)
	s.Assert().Equal(16, s.abilities.Str)
	s.Assert().Equal(15, s.abilities.Int)
	s.Assert().Equal(0, s.abilities.Wis)
	s.Assert().Equal(0, s.abilities.Cha)

	s.abilities = s.abilities.fillEmptyScoreByPriority(14)
	s.Assert().Equal(18, s.abilities.Con)
	s.Assert().Equal(17, s.abilities.Dex)
	s.Assert().Equal(16, s.abilities.Str)
	s.Assert().Equal(15, s.abilities.Int)
	s.Assert().Equal(14, s.abilities.Wis)
	s.Assert().Equal(0, s.abilities.Cha)

	s.abilities = s.abilities.fillEmptyScoreByPriority(13)
	s.Assert().Equal(18, s.abilities.Con)
	s.Assert().Equal(17, s.abilities.Dex)
	s.Assert().Equal(16, s.abilities.Str)
	s.Assert().Equal(15, s.abilities.Int)
	s.Assert().Equal(14, s.abilities.Wis)
	s.Assert().Equal(13, s.abilities.Cha)

	s.abilities = s.abilities.fillEmptyScoreByPriority(12)
	s.Assert().Equal(18, s.abilities.Con)
	s.Assert().Equal(17, s.abilities.Dex)
	s.Assert().Equal(16, s.abilities.Str)
	s.Assert().Equal(15, s.abilities.Int)
	s.Assert().Equal(14, s.abilities.Wis)
	s.Assert().Equal(13, s.abilities.Cha)
}
