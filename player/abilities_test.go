package player

import (
	"github.com/stretchr/testify/suite"
	"testing"
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

func (s *AbilitiesTestSuite) Test_Modifiers() {
	s.Assert().Equal(-5, AbilityScoreModifier(1))
	s.Assert().Equal(-5, AbilityScoreModifier(0))
	s.Assert().Equal(-5, AbilityScoreModifier(-1))

	s.Assert().Equal(10, AbilityScoreModifier(35))

	s.Assert().Equal(4, AbilityScoreModifier(18))
	s.Assert().Equal(0, AbilityScoreModifier(11))
}

func (s *AbilitiesTestSuite) Test_SetScore() {
	s.abilities.SetScore("str", 18)
	s.Assert().Equal(AbilityScore(18), s.abilities.Strength)
	s.Assert().Equal(AbilityScore(0), s.abilities.Intelligence)
	s.Assert().Equal(AbilityScore(0), s.abilities.Dexterity)
	s.Assert().Equal(AbilityScore(0), s.abilities.Wisdom)
	s.Assert().Equal(AbilityScore(0), s.abilities.Constitution)
	s.Assert().Equal(AbilityScore(0), s.abilities.Charisma)
}
