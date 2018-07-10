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
	fighter := NewTestCombatant("fighter")
	fightee := NewTestCombatant("fightee")

	log.Printf("%p", fighter)
	log.Printf("%p", fightee)

	suite.Assert().False(suite.fightLedger.IsFighting(fighter))
	suite.Assert().False(suite.fightLedger.IsFighting(fightee))

	suite.fightLedger.Fight(fighter, fightee, "zoneId", "roomId")

	// fighting is one directional
	suite.Assert().True(suite.fightLedger.IsFighting(fighter))
	suite.Assert().False(suite.fightLedger.IsFighting(fightee))
}

func (suite *FightLedgerSuite) TestAlreadyFighting() {
	fighter := NewTestCombatant("fighter")
	fightee := NewTestCombatant("fightee")

	suite.Assert().NoError(suite.fightLedger.Fight(fighter, fightee, "zoneId", "roomId"))
	// can't fight when you're already fighting
	suite.Assert().Error(suite.fightLedger.Fight(fighter, fightee, "zoneId", "roomId"))

}
