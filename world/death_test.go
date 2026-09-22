package world

import (
	"github.com/trasa/watchmud/event"
)

// A player who dies leaves no corpse: they wake up in the death room with a
// single hit point. These ride on durabilitySuite for its dies() helper.

func (s *durabilitySuite) TestDyingLeavesYouWithOneHitPoint() {
	s.dies()

	s.Assert().False(s.p.Dead())
	s.Assert().Equal(1, s.p.CurrentHealth())
}

// Pointed somewhere other than where the fight was, so the move is visible.
func (s *durabilitySuite) TestDyingTakesYouToTheDeathRoom() {
	s.w.DeathRoom = s.w.VoidRoom

	s.dies()

	s.Assert().Equal(s.w.VoidRoom, s.w.getPlayerRoom(s.p))
	s.Assert().Empty(s.w.StartRoom.Players(), "gone from where they died")
	s.Assert().Len(s.w.VoidRoom.Players(), 1)

	desc := sent[event.RoomDescription](s.T(), s.r, len(s.r.Sent)-1)
	s.Assert().Equal(s.w.VoidRoom.Name, desc.Name, "and shown where they woke up")
}

func (s *durabilitySuite) TestDyingLeavesNoCorpse() {
	before := s.w.StartRoom.Inventory.Len()

	s.dies()

	s.Assert().Equal(before, s.w.StartRoom.Inventory.Len())
}

func (s *durabilitySuite) TestDyingEndsTheFight() {
	s.dies()

	s.Assert().False(s.w.fightLedger.InFight(s.p))
	s.Assert().False(s.w.fightLedger.InFight(s.drone))
}

func (s *durabilitySuite) TestTheRoomIsToldAPlayerDied() {
	s.dies()

	d := sent[event.Died](s.T(), s.r, 1)
	s.Assert().Equal("victim", d.Target)
	s.Assert().True(d.IsPlayer)
}

// Death happens on a pulse, with no command after it to save the player, so
// it saves them itself: a crash straight after must not undo it.
func (s *durabilitySuite) TestDyingIsSaved() {
	s.w.DeathRoom = s.w.VoidRoom

	s.dies()

	rec, found, err := s.w.store.Load(s.p.Name())
	s.Require().NoError(err)
	s.Require().True(found)
	s.Assert().Equal(1, rec.CurHealth)
	s.Assert().Equal(s.w.VoidRoom.Zone.Id, rec.LastZoneId)
	s.Assert().Equal(s.w.VoidRoom.Id, rec.LastRoomId)
}
