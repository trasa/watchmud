package combat

import (
	"github.com/stretchr/testify/suite"
	"math/rand"
	"testing"
)

type MeleeSuite struct {
	suite.Suite
	fighter Combatant
	victim  Combatant
}

func TestMeleeSuite(t *testing.T) {
	suite.Run(t, new(MeleeSuite))
}

func (suite *MeleeSuite) SetupTest() {
	suite.fighter = NewTestCombatant("fighter")
	suite.victim = NewTestCombatant("victim")
}

func (suite *MeleeSuite) TestHitFailed() {
	r := rand.New(rand.NewSource(1))
	noChance := meleeAttack(suite.fighter, suite.victim, r, 0.00)
	suite.Assert().False(noChance.WasHit)
	suite.Assert().Equal(0, noChance.Damage)
}

func (suite *MeleeSuite) TestHitSuccess() {
	r := rand.New(rand.NewSource(1))
	noChance := meleeAttack(suite.fighter, suite.victim, r, 1.00)
	suite.Assert().True(noChance.WasHit)
	suite.Assert().Equal(5, noChance.Damage)
}
