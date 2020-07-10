package combat

import (
	"github.com/justinian/dice"
	"github.com/stretchr/testify/suite"
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
	suite.fighter = NewTestCombatant("fighter", 10, []DamageType{}, []DamageType{})
	suite.victim = NewTestCombatant("victim", 10, []DamageType{}, []DamageType{})
}

func (suite *MeleeSuite) TestHitFailed() {
	rollResult := dice.StdResult{Total: 2}

	noChance := meleeAttack(suite.fighter, suite.victim, rollResult)
	suite.Assert().False(noChance.WasHit)
	suite.Assert().Equal(int64(0), noChance.Damage)
}

func (suite *MeleeSuite) TestHitSuccess() {
	rollResult := dice.StdResult{Total: 16}

	noChance := meleeAttack(suite.fighter, suite.victim, rollResult)
	suite.Assert().True(noChance.WasHit)
	// TODO calculate damage
}

// TODO critical success
// TODO critical fail
// TODO resistance
// TODO vulnerability
