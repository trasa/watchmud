package combat

import (
	"github.com/stretchr/testify/suite"
	"log"
	"testing"
)

type FightLedgerSuite struct {
	suite.Suite
	fightLedger *FightLedger
}

func TestFightLedgerSuite(t *testing.T) {
	suite.Run(t, new(FightLedgerSuite))
}

func (suite *FightLedgerSuite) SetupTest() {
	suite.fightLedger = NewFightLedger()
}

func (suite *FightLedgerSuite) TestIsFighting() {
	fighter := NewTestCombatant()
	fightee := NewTestCombatant()

	log.Printf("%p", fighter)
	log.Printf("%p", fightee)

	suite.Assert().False(suite.fightLedger.IsFighting(fighter))
	suite.Assert().False(suite.fightLedger.IsFighting(fightee))

	suite.fightLedger.Fight(fighter, fightee)

	// fighting is one directional
	suite.Assert().True(suite.fightLedger.IsFighting(fighter))
	suite.Assert().False(suite.fightLedger.IsFighting(fightee))
}

func (suite *FightLedgerSuite) TestAlreadyFighting() {
	fighter := NewTestCombatant()
	fightee := NewTestCombatant()

	suite.Assert().NoError(suite.fightLedger.Fight(fighter, fightee))
	// can't fight when you're already fighting
	suite.Assert().Error(suite.fightLedger.Fight(fighter, fightee))

}
