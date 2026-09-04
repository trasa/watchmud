package player

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"github.com/trasa/watchmud-message/slot"
	"github.com/trasa/watchmud/object"
)

type PlayerSuite struct {
	suite.Suite
	p *Player
	r *Recorder
}

func TestPlayerSuite(t *testing.T) {
	suite.Run(t, new(PlayerSuite))
}

func (s *PlayerSuite) SetupTest() {
	s.r = &Recorder{}
	s.p = NewTestPlayer("testdood", "testdood", s.r)
}

func (s *PlayerSuite) TestAddInventory_New() {
	// old test, do not trust it

	defnPtr := object.NewDefinition("defnid", "name", "zone",
		object.Food, []string{}, "short desc", "in room", slot.None)
	instPtr := &object.Instance{
		InstanceId: uuid.New(),
		Definition: defnPtr,
	}

	s.p.inventory.Add(instPtr)

	invs := s.p.inventory.GetAll()
	s.Assert().Equal(1, len(invs))
	obj := invs[0]
	s.Assert().Equal(instPtr.Id(), obj.Id())
	s.Assert().Equal("defnid", obj.Definition.Identifier())
}

func (s *PlayerSuite) TestMeleeDamage() {
	// old test, do not trust it

	startingHealth := s.p.curHealth

	isDead := s.p.TakeMeleeDamage(5)

	s.Assert().False(isDead)
	s.Assert().Equal(startingHealth-5, s.p.curHealth)
}

func (s *PlayerSuite) TestFatalMeleeDamage() {
	// old test, do not trust it

	startingHealth := s.p.curHealth

	isDead := s.p.TakeMeleeDamage(startingHealth)

	s.Assert().True(isDead)
	s.Assert().Equal(int64(0), s.p.curHealth)
}

func (s *PlayerSuite) TestOverwhelminglyFatalMeleeDamage() {
	// old test, do not trust it

	startingHealth := s.p.curHealth

	isDead := s.p.TakeMeleeDamage(startingHealth * 2)

	s.Assert().True(isDead)
	s.Assert().True(s.p.curHealth < 0)
}

func (s *PlayerSuite) TestIsDead() {
	// old test, do not trust it

	s.Assert().False(s.p.IsDead())
	s.p.curHealth = 0
	s.Assert().True(s.p.IsDead())
}
