package memstore

import (
	"testing"
	"uuid"

	"github.com/stretchr/testify/suite"
	"github.com/trasa/watchmud/player"
)

type storeSuite struct {
	suite.Suite
	store  *Store
	player *player.Player
}

func TestStoreSuite(t *testing.T) {
	suite.Run(t, new(storeSuite))
}
func (s *storeSuite) SetupTest() {
	s.store = New()
	s.player = player.NewTestPlayer(
		uuid.New(),
		"test",
		&player.Recorder{},
	)
}

func (s *storeSuite) TestSavePlayer() {
	rec := s.player.Record()
	s.Assert().Equal("passwordHash", rec.PasswordHash)
	s.Require().NoError(s.store.Save(rec))

	retrieved, found, err := s.store.Load(s.player.Name())
	s.Assert().True(found)
	s.Assert().NoError(err)
	s.Assert().Equal(retrieved.PasswordHash, rec.PasswordHash)
}
