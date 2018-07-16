package mobile

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type InstanceMeleeDamageSuite struct {
	suite.Suite
	instance *Instance
}

func TestInstanceMeleeDamageSuite(t *testing.T) {
	suite.Run(t, new(InstanceMeleeDamageSuite))
}

func (s *InstanceMeleeDamageSuite) SetupTest() {
	defn := NewDefinition("id", "damaged", "testZone", []string{}, "", "",
		25,
		WanderingDefinition{
			CanWander: false,
		})
	defn.MaxHealth = 100
	s.instance = NewInstance(defn)
	s.instance.CurHealth = 100
}

func (s *InstanceMeleeDamageSuite) TestMeleeDamage() {
	startingHealth := s.instance.CurHealth

	isDead := s.instance.TakeMeleeDamage(5)

	s.Assert().False(isDead)
	s.Assert().Equal(startingHealth-5, s.instance.CurHealth)
}

func (s *InstanceMeleeDamageSuite) TestFatalMeleeDamage() {

	isDead := s.instance.TakeMeleeDamage(s.instance.CurHealth)

	s.Assert().True(isDead)
	s.Assert().Equal(0, s.instance.CurHealth)
}

func (s *InstanceMeleeDamageSuite) TestOverwhelmingFatalMeleeDamage() {

	isDead := s.instance.TakeMeleeDamage(s.instance.CurHealth * 2)

	s.Assert().True(isDead)
	s.Assert().True(s.instance.CurHealth < 0)
}

func (s *InstanceMeleeDamageSuite) TestIsDead() {
	s.Assert().False(s.instance.IsDead())
	s.instance.CurHealth = 0
	s.Assert().True(s.instance.IsDead())
}
