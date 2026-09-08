package memstore

import (
	"fmt"
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

func (s *Store) Create(r *player.Record) (*player.Record, error) {
	if _, found := s.byName[r.Name]; found {
		return nil, fmt.Errorf("create: record with this name already exists: %s", r.Name)
	}
	if _, found := s.byId[r.Id]; found {
		return nil, fmt.Errorf("create: record with this id already exists: %s", r.Id)
	}
	s.byName[r.Name] = r
	s.byId[r.Id] = r
	return r, nil
}

func (s *Store) Save(r *player.Record) error {
	if _, found := s.byName[r.Name]; !found {
		return fmt.Errorf("save: record with this name does not exist: %s", r.Name)
	}
	if _, found := s.byId[r.Id]; !found {
		return fmt.Errorf("save: record with this id does not exist: %s", r.Id)
	}
	s.byName[r.Name] = r
	s.byId[r.Id] = r
	return nil
}
