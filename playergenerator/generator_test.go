package playergenerator

import (
	"github.com/stretchr/testify/suite"
	"github.com/trasa/watchmud/db"
	"github.com/trasa/watchmud/player"
	"testing"
)

type GeneratorTestSuite struct {
	suite.Suite
}

func TestGeneratorTestSuite(t *testing.T) {
	suite.Run(t, new(GeneratorTestSuite))
}

func (s *GeneratorTestSuite) SetupTest() {
}

func (s *GeneratorTestSuite) TestAbilities() {
	race := db.RaceData{
		StrBonus: 1,
	}
	class := db.ClassData{
		AbilityPreference: &db.AbilityPreferenceList{
			Preferences: []string{"str", "dex", "con"},
		},
	}
	a := generateAbilities(race, class)

	s.Assert().Equal(player.AbilityScore(16), a.Strength)
	s.Assert().Equal(player.AbilityScore(14), a.Dexterity)
	s.Assert().Equal(player.AbilityScore(13), a.Constitution)
	s.Assert().Equal(player.AbilityScore(12), a.Intelligence)
	s.Assert().Equal(player.AbilityScore(10), a.Wisdom)
	s.Assert().Equal(player.AbilityScore(8), a.Charisma)
}
