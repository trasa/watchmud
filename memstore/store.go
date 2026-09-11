package memstore

import (
	"uuid"

	"github.com/trasa/watchmud/player"
)

type Store struct {
	byId   map[uuid.UUID]*player.Record
	byName map[string]*player.Record
}

func New() *Store {
	return &Store{
		byId:   make(map[uuid.UUID]*player.Record),
		byName: make(map[string]*player.Record),
	}
}

func (s *Store) Load(name string) (*player.Record, bool, error) {
	r, found := s.byName[name]
	return r, found, nil
}

func (s *Store) Save(r *player.Record) error {
	s.byName[r.Name] = r
	s.byId[r.Id] = r
	return nil
}
