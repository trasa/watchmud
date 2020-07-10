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
	fighter := NewTestCombatant("fighter", 10, []DamageType{}, []DamageType{})
	fightee := NewTestCombatant("fightee", 10, []DamageType{}, []DamageType{})

	log.Printf("%p", fighter)
	log.Printf("%p", fightee)

	suite.Assert().False(suite.fightLedger.IsFighting(fighter))
	suite.Assert().False(suite.fightLedger.IsFighting(fightee))

	suite.fightLedger.Fight(fighter, fightee, "zoneId", "roomId")

	// fighting starts a bidirectional fight
	suite.Assert().True(suite.fightLedger.IsFighting(fighter))
	suite.Assert().True(suite.fightLedger.IsFighting(fightee))
}

func (suite *FightLedgerSuite) TestAlreadyFighting() {
	fighter := NewTestCombatant("fighter", 10, []DamageType{}, []DamageType{})
	fightee := NewTestCombatant("fightee", 10, []DamageType{}, []DamageType{})

	suite.Assert().NoError(suite.fightLedger.Fight(fighter, fightee, "zoneId", "roomId"))
	// can't fight when you're already fighting
	suite.Assert().Error(suite.fightLedger.Fight(fighter, fightee, "zoneId", "roomId"))
}

func (suite *FightLedgerSuite) TestFightingSomeoneWhoIsFighting() {
	fighter := NewTestCombatant("fighter", 10, []DamageType{}, []DamageType{})
	otherFighter := NewTestCombatant("otherFighter", 10, []DamageType{}, []DamageType{})
	fightee := NewTestCombatant("fightee", 10, []DamageType{}, []DamageType{})

	suite.Assert().NoError(suite.fightLedger.Fight(otherFighter, fightee, "zoneId", "roomId"))
	// can be fought by more than 1
	suite.Assert().NoError(suite.fightLedger.Fight(fighter, fightee, "zoneId", "roomId"))
	suite.Assert().Equal(otherFighter, suite.fightLedger.fightMap[fightee].Fightee)
}

func (suite *FightLedgerSuite) TestEndFight() {
	fighter := NewTestCombatant("fighter", 10, []DamageType{}, []DamageType{})
	fightee := NewTestCombatant("fightee", 10, []DamageType{}, []DamageType{})

	suite.Assert().NoError(suite.fightLedger.Fight(fighter, fightee, "zoneId", "roomId"))

	suite.fightLedger.EndFight(fighter)
	suite.Assert().False(suite.fightLedger.IsFighting(fighter))
	suite.Assert().True(suite.fightLedger.IsFighting(fightee))
}
