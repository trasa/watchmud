package object

type Category int

const (
	WEAPON Category = iota
	WAND
	STAFF
	TREASURE
	ARMOR
	FOOD
	OTHER
)

// Definition of what it means to be an "object"
// the "platonic form" of an object, if you will.
//
// There is a definition of the type "ShinySword" which defines the
// properties of what it means to be a ShinySword. An actual
// ShinySword lying around is an Instance.
type Definition interface {
	TypeId() string
	Aliases() []string
	Name() string
	ShortDescription() string    // description of the object when being used: "a long, green stick" -> "The Beastly Fido picks up the long, green stick."
	DescriptionOnGround() string // description of the object when lying on the ground: "A shiny sword is lying here."
}

// The Instances of the Definitions in the world around you.
// That ShinySword in your hand has certain properties, some of
// which were inherited by what it means to be a ShinySword (ObjectType)
// and others which have happened to that particular instance
// (soul-bound to you, made invisible, with some damage to the hilt).
type Instance interface {
	TypeId() string // refers to the Definition from which this Instance came
	Name() string
}
