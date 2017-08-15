package object

// The Instances of the Definitions in the world around you.
// That ShinySword in your hand has certain properties, some of
// which were inherited by what it means to be a ShinySword (ObjectType)
// and others which have happened to that particular instance
// (soul-bound to you, made invisible, with some damage to the hilt).
type Instance struct {
	// TODO what properties does an Instance have?
	Definition *Definition
}

func NewInstance(defn *Definition) *Instance {
	return &Instance{
		Definition: defn,
	}
}
