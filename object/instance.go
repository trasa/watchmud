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

	// Durability is what this particular one has left, counting down from
	// Definition.MaxDurability. See durability.go.
	Durability int
}

// IdStr from the Thing interface
// TODO figure out if this can be removed
func (i *Instance) IdStr() string { return i.Id.String() }

func (i *Instance) IsGettable() bool {
	return i.Definition.Gettable()
}

func (i *Instance) Matches(target string) bool {
	return i.Definition.Matches(target)
}

// NewInstance of a definition, brand new: full durability. Anything restoring
// one that has already been used -- a save file -- sets Durability afterwards.
func NewInstance(id uuid.UUID, d *Definition) *Instance {
	return &Instance{
		Id:         id,
		Definition: d,
		Durability: d.MaxDurability,
	}
}
