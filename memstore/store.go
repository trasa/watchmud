package memstore

import (
	"errors"

	"github.com/trasa/watchmud/player"
)

type Store struct {
	byName map[string]*player.Record
}

func New() *Store {
	return &Store{
		byName: make(map[string]*player.Record),
	}
}

func (s *Store) Load(name string) (*player.Record, bool, error) {
	r, found := s.byName[name]
	return r, found, nil
}

func (s *Store) Create(r *player.Record) (*player.Record, error) {
	if _, found := s.byName[r.Name]; found {
		return nil, errors.New("record with this name already exists")
	}
	s.byName[r.Name] = r
	return r, nil
}

func (s *Store) Save(r *player.Record) error {
	if _, found := s.byName[r.Name]; !found {
		return errors.New("record with this name does not exist")
	}
	s.byName[r.Name] = r
	return nil
}
