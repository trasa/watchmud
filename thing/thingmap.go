package thing

import (
	"errors"
	"fmt"
)

//noinspection GoNameStartsWithPackageName
type Thing interface {
	Id() string
}

// map InstanceIDs to Instance objects
type Map map[string]Thing

// Add an instance to the map.
// If it's already there, thats' an error.
func (m Map) Add(t Thing) error {
	if _, exists := m[t.Id()]; exists {
		return errors.New(fmt.Sprintf("instance %s already exists in map", t.Id()))
	} else {
		m[t.Id()] = t
	}
	return nil
}

// Remove an instance from the map.
// If it doesn't exist in the map, that's an error.
func (m Map) Remove(t Thing) error {
	if _, exists := m[t.Id()]; !exists {
		return errors.New(fmt.Sprintf("instance %s can't be removed, it doesn't exist in map", t.Id()))
	} else {
		delete(m, t.Id())
	}
	return nil
}
