package spaces

import (
	"fmt"
	"iter"
	"uuid"

	"github.com/watchmud/watchmud/object"
	"github.com/watchmud/watchmud/ordered"
)

// RoomInventory is what is lying on the floor of a room, in the order it
// landed there.
type RoomInventory struct {
	objects *ordered.List[uuid.UUID, *object.Instance]
}

func NewRoomInventory() *RoomInventory {
	return &RoomInventory{
		objects: ordered.NewList(func(o *object.Instance) uuid.UUID { return o.Id }),
	}
}

// All the instances in the room, in the order they were dropped.
func (ri *RoomInventory) All() iter.Seq[*object.Instance] {
	return ri.objects.All()
}

// Len is the number of instances in the room.
func (ri *RoomInventory) Len() int {
	return ri.objects.Len()
}

// InstanceId finds the instance with this id in the room
func (ri *RoomInventory) InstanceId(id uuid.UUID) (*object.Instance, bool) {
	return ri.objects.Get(id)
}

// NameOrAlias finds the instances this target names, in the order they were
// dropped.
// TODO this needs to become much more sophisticated...
// Note that there is much left undone by this implementation
// (stacks of items, aliases...)
func (ri *RoomInventory) NameOrAlias(target string) []*object.Instance {
	return ordered.FindAll(ri.objects, target)
}

func (ri *RoomInventory) Add(inst *object.Instance) error {
	if err := ri.objects.Add(inst); err != nil {
		return fmt.Errorf("add object %s to room inventory: %w", inst.Id, err)
	}
	return nil
}

func (ri *RoomInventory) Remove(inst *object.Instance) error {
	if err := ri.objects.Remove(inst); err != nil {
		return fmt.Errorf("remove object %s from room inventory: %w", inst.Id, err)
	}
	return nil
}
