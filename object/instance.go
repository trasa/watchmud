package object

import "uuid"

// Instance of the Definitions in the world around you.
// That ShinySword in your hand has certain properties, some of
// which were inherited by what it means to be a ShinySword (Definition)
// and others which have happened to that particular instance
// (soul-bound to you, made invisible, with some damage to the hilt).
type Instance struct {
	Id         uuid.UUID
	Definition *Definition
}

// IdStr from the Thing interface TODO figure out if this can be removed
func (i *Instance) IdStr() string { return i.Id.String() }

func (i *Instance) CanEquipWeapon() bool {
	return i.Definition.CanEquipWeapon()
}

func (i *Instance) IsGettable() bool {
	return i.Definition.IsGettable()
}

func NewInstance(d *Definition) *Instance {
	return NewInstanceWithId(uuid.New(), d)
}

func NewInstanceWithId(id uuid.UUID, d *Definition) *Instance {
	return &Instance{
		Id:         id,
		Definition: d,
	}
}
