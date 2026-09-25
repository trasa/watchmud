package world

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/watchmud/watchmud/event"
)

// What happens to everyone else's fights when one combatant dies. These drive
// the death handling directly, since that is where the bookkeeping lives --
// DoViolence itself rolls through w.roller and can be driven a swing at a
// time, which violence_armorclass_test.go does.
type violenceSuite struct {
	worldTestSuite
}

func TestViolenceSuite(t *testing.T) {
	suite.Run(t, new(violenceSuite))
}

// Two drones on one player. Killing one must leave the other one fighting --
// this ended every fight the winner was in, so a second attacker silently
// dropped out the moment the first one died.
func (s *violenceSuite) TestKillingOneAttackerLeavesTheOtherFighting() {
	target, exists := s.w.StartRoom.FindMobile("target")
	s.Require().True(exists)
	little, exists := s.w.StartRoom.FindMobile("little")
	s.Require().True(exists)

	s.Require().NoError(s.w.fightLedger.Fight(target, s.p, "wrathrock", "temple_square"))
	s.Require().NoError(s.w.fightLedger.Fight(little, s.p, "wrathrock", "temple_square"))
	s.Require().True(s.w.fightLedger.IsFighting(little))

	s.w.combatantDied(target, s.w.StartRoom, true)

	s.Assert().False(s.w.fightLedger.IsFighting(target), "the dead one stops fighting")
	s.Assert().True(s.w.fightLedger.IsFighting(little), "the survivor keeps fighting")
	s.Assert().True(s.w.fightLedger.InFight(s.p), "and the player is still in a fight")
}

// Your target dies while something else is still hitting you: you turn on
// it, rather than stand there "in a fight" and never swing again. The same
// rule is what turns the Barrow-King on the next player once the tank falls.
func (s *violenceSuite) TestWhenYourTargetDiesYouTurnOnTheNextAttacker() {
	target, exists := s.w.StartRoom.FindMobile("target")
	s.Require().True(exists)
	little, exists := s.w.StartRoom.FindMobile("little")
	s.Require().True(exists)

	s.Require().NoError(s.w.fightLedger.Fight(s.p, target, "wrathrock", "temple_square"))
	s.Require().NoError(s.w.fightLedger.Fight(little, s.p, "wrathrock", "temple_square"))

	s.w.combatantDied(target, s.w.StartRoom, true)

	fight := s.w.fightLedger.GetFight(s.p)
	s.Require().NotNil(fight, "the player is fighting again")
	s.Assert().Same(little, fight.Fightee)
}

// The one who died leaves in both directions: nothing is left swinging at a
// corpse either.
func (s *violenceSuite) TestTheDeadLeaveInBothDirections() {
	target, exists := s.w.StartRoom.FindMobile("target")
	s.Require().True(exists)

	s.Require().NoError(s.w.fightLedger.Fight(s.p, target, "wrathrock", "temple_square"))

	s.w.combatantDied(target, s.w.StartRoom, true)

	s.Assert().False(s.w.fightLedger.InFight(target))
	s.Assert().False(s.w.fightLedger.InFight(s.p))
}

func (s *violenceSuite) TestTheRoomIsTold() {
	target, exists := s.w.StartRoom.FindMobile("target")
	s.Require().True(exists)

	s.w.combatantDied(target, s.w.StartRoom, true)

	s.Assert().Equal("Target Drone", sent[event.Died](s.T(), s.r, 0).Target)
}

// A fight can outlive its room (the ledger snapshots zone and room ids when
// it starts), and then there is nobody to tell.
func (s *violenceSuite) TestNoRoomNoNotification() {
	target, exists := s.w.StartRoom.FindMobile("target")
	s.Require().True(exists)

	s.w.combatantDied(target, nil, false)

	s.Assert().Empty(s.r.Sent)
}
