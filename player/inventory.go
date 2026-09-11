package player

import (
	"maps"
	"slices"
	"uuid"

	"github.com/trasa/watchmud/object"
)

type Inventory struct {
	byId map[uuid.UUID]*object.Instance // instance_id -> instance obj
}

func NewInventory() *Inventory {
	return &Inventory{
		byId: make(map[uuid.UUID]*object.Instance),
	}
}

func (pi *Inventory) GetAll() []*object.Instance {
	v := slices.Collect(maps.Values(pi.byId))
	// v.sort() // TODO
	return v
}

/*
func (pi *Inventory) sort() {
	sort.SliceStable(pi.sorted, func(i, j int) bool {
		return pi.sorted[i].InstanceId.String() < pi.sorted[j].InstanceId.String()
	})
}
*/

func (pi *Inventory) ByInstanceId(id uuid.UUID) (*object.Instance, bool) {
	inst, exists := pi.byId[id]
	return inst, exists
}

// GetByNameOrAlias the items that match this string
func (pi *Inventory) GetByNameOrAlias(target string) (objects []*object.Instance) {
	// TODO handle case where target is "2.knife" (return the 2nd knife)
	// TODO handle all the other target cases
	objects = []*object.Instance{}
	for _, obj := range pi.GetAll() {
		if obj.Definition.Name == target || obj.Definition.HasAlias(target) {
			objects = append(objects, obj)
		}
	}
	return objects
}

// Add an object into the inventory
func (pi *Inventory) Add(inst *object.Instance) {
	pi.byId[inst.Id] = inst
}

// Remove an object from the inventory
func (pi *Inventory) Remove(inst *object.Instance) {
	delete(pi.byId, inst.Id)
}
