package rules

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type RolesTestSuite struct {
	suite.Suite
	cat *Catalog
}

func TestRolesTestSuite(t *testing.T) {
	suite.Run(t, new(RolesTestSuite))
}

func (s *RolesTestSuite) SetupTest() {
	cat, err := NewTestCatalog()
	require.NoError(s.T(), err)
	s.cat = cat
}

func (s *RolesTestSuite) Test_NoGearIsNoRole() {
	s.Assert().Nil(s.cat.RoleFor(nil))
	s.Assert().Nil(s.cat.RoleFor(map[string]int{}))
}

func (s *RolesTestSuite) Test_GearThatContributesNothingIsNoRole() {
	s.Assert().Nil(s.cat.RoleFor(map[string]int{"tank": 0, "healer": 0}))
}

func (s *RolesTestSuite) Test_HighestTotalWins() {
	r := s.cat.RoleFor(map[string]int{"tank": 2, "healer": 5, "striker": 1})
	s.Require().NotNil(r)
	s.Assert().Equal("healer", r.Id)
}

func (s *RolesTestSuite) Test_TieGoesToTheFirstDeclaredRole() {
	// tank is declared before striker in NewTestRoles
	r := s.cat.RoleFor(map[string]int{"striker": 4, "tank": 4})
	s.Require().NotNil(r)
	s.Assert().Equal("tank", r.Id)
}

func (s *RolesTestSuite) Test_UnknownRoleIdsAreIgnored() {
	// content validation rejects these at load time; RoleFor must not be
	// fooled by one that somehow gets this far.
	s.Assert().Nil(s.cat.RoleFor(map[string]int{"bard": 99}))
}

func (s *RolesTestSuite) Test_DuplicateRoleIdIsAnError() {
	_, err := NewCatalog(NewTestSpecies(), []*Role{
		{Id: "tank", Name: "Tank"},
		{Id: "tank", Name: "Also Tank"},
	})
	s.Assert().ErrorContains(err, "duplicate role id")
}

func (s *RolesTestSuite) Test_MissingRoleIdIsAnError() {
	_, err := NewCatalog(NewTestSpecies(), []*Role{{Name: "Nameless"}})
	s.Assert().ErrorContains(err, "missing id")
}

func (s *RolesTestSuite) Test_DefaultLineage() {
	l := s.cat.DefaultLineage()
	s.Require().NotNil(l)
	s.Assert().Equal("human", l.Id)
}
