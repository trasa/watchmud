package combat

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/trasa/watchmud/testdice"
)

type MeleeSuite struct {
	suite.Suite
	fighter Combatant
	victim  Combatant
	roller  *testdice.LoadedDice
}

func TestMeleeSuite(t *testing.T) {
	suite.Run(t, new(MeleeSuite))
}

func (s *MeleeSuite) SetupTest() {
	s.fighter = NewTestCombatant("fighter", 10, []DamageType{}, []DamageType{})
	s.victim = NewTestCombatant("victim", 10, []DamageType{}, []DamageType{})
	s.roller = testdice.New()
}

func (s *MeleeSuite) TestHitFailed() {
	s.roller.Add(2)

	noChance, err := AttemptMeleeAttack(s.roller, s.fighter, s.victim)

	s.Assert().NoError(err)
	s.Assert().False(noChance.WasHit)
	s.Assert().Equal(0, noChance.Damage)
}

func (s *MeleeSuite) TestHitSuccess() {
	s.roller.Load([]int{19, 3})
	result, err := AttemptMeleeAttack(s.roller, s.fighter, s.victim)
	s.Assert().NoError(err)
	s.Assert().True(result.WasHit)
	s.Assert().Equal(3, result.Damage)
}

// TODO critical success
// TODO critical fail
// TODO resistance
// TODO vulnerability
