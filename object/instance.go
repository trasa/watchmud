package object

import "github.com/satori/go.uuid"

// The Instances of the Definitions in the world around you.
// That ShinySword in your hand has certain properties, some of
// which were inherited by what it means to be a ShinySword (Definition)
// and others which have happened to that particular instance
// (soul-bound to you, made invisible, with some damage to the hilt).
type Instance struct {
	InstanceId uuid.UUID
	Definition *Definition
}

func (i *Instance) Id() string {
	return i.InstanceId.String()
}

func NewInstance(defn *Definition) *Instance {
	return &Instance{
		InstanceId: uuid.NewV4(),
		Definition: defn,
	}
}
