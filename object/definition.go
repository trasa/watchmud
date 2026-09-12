package object

import (
	"slices"

	"github.com/trasa/watchmud/behavior"
	"github.com/trasa/watchmud/slot"
)

// Definition of what it means to be an "object"
// the "platonic form" of an object, if you will.
//
// There is a definition of the type "ShinySword" which defines the
// properties of what it means to be a ShinySword. An actual
// ShinySword lying around is an Instance.
type Definition struct {
	ObjectId            Id
	Aliases             []string
	Categories          CategorySet
	Name                string
	ShortDescription    string // description of the object when being used: "a long, green stick" -> "The Beastly Fido picks up the long, green stick."
	DescriptionOnGround string // description of the object when lying on the ground: "A shiny sword is lying here."
	WearLocation        slot.Location
	Behaviors           behavior.BehaviorSet

	// RoleWeights is what this object contributes to each role (by
	// rules.Role.Id) while it is equipped: {"tank": 3} is solidly tanky gear,
	// {"healer": 2, "tank": 1} is a healer's kit with some heft to it. The
	// weights of everything a character has equipped are summed and the
	// highest total is their role -- so this field, across a few objects, is
	// the entire replacement for character classes.
	//
	// Set by the loader after construction, like Behaviors, and validated
	// against the catalog there: an unknown role id is a content error.
	// Nil for anything that isn't equipment, which is most objects.
	RoleWeights map[string]int
}

func NewDefinition(
	id string,
	name string,
	zoneId string,
	category Category,
	aliases []string,
	shortDescription string,
	descriptionOnGround string,
	wearLocation slot.Location) *Definition {
	d := &Definition{
		ObjectId:            NewObjectId(id, zoneId),
		Name:                name,
		ShortDescription:    shortDescription,
		DescriptionOnGround: descriptionOnGround,
		Categories:          make(CategorySet),
		Aliases:             aliases,
		WearLocation:        wearLocation,
		Behaviors:           behavior.NewBehaviorSet(),
	}
	d.Categories.Add(category)
	return d
}

func (d *Definition) IsWeapon() bool {
	return d.Categories.Contains(Weapon)
}

func (d *Definition) NoTake() bool {
	return d.Behaviors.Contains(behavior.NoTake)
}

func (d *Definition) Gettable() bool {
	return !d.NoTake() // there might be other reasons why you can't get it in the future
}

func (d *Definition) Wearable() bool {
	return d.WearLocation != slot.None
}

func (d *Definition) HasAlias(target string) bool {
	return slices.Contains(d.Aliases, target)
}
