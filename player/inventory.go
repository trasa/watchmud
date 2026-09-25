package player

import (
	"iter"
	"uuid"

	"github.com/rs/zerolog/log"
	"github.com/watchmud/watchmud/object"
	"github.com/watchmud/watchmud/ordered"
)

// Inventory is what a player is carrying, in the order they picked it up.
type Inventory struct {
	objects *ordered.List[uuid.UUID, *object.Instance]
}

func NewInventory() *Inventory {
	return &Inventory{
		objects: ordered.NewList(func(o *object.Instance) uuid.UUID { return o.Id }),
	}
}

// All the items being carried, in the order they were picked up.
func (pi *Inventory) All() iter.Seq[*object.Instance] {
	return pi.objects.All()
}

// Len is the number of items being carried.
func (pi *Inventory) Len() int {
	return pi.objects.Len()
}

func (pi *Inventory) ByInstanceId(id uuid.UUID) (*object.Instance, bool) {
	return pi.objects.Get(id)
}

// GetByNameOrAlias the items that match this string, oldest first
// TODO handle case where target is "2.knife" (return the 2nd knife)
// TODO handle all the other target cases
func (pi *Inventory) GetByNameOrAlias(target string) []*object.Instance {
	return ordered.FindAll(pi.objects, target)
}

// Add an object into the inventory
func (pi *Inventory) Add(inst *object.Instance) {
	if err := pi.objects.Add(inst); err != nil {
		log.Warn().Err(err).Str("instance", inst.IdStr()).Msg("Inventory.Add")
	}
}

// Remove an object from the inventory
func (pi *Inventory) Remove(inst *object.Instance) {
	if err := pi.objects.Remove(inst); err != nil {
		log.Warn().Err(err).Str("instance", inst.IdStr()).Msg("Inventory.Remove")
	}
}
