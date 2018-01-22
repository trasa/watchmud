package object

import "github.com/trasa/watchmud/slot"

// Definition of what it means to be an "object"
// the "platonic form" of an object, if you will.
//
// There is a definition of the type "ShinySword" which defines the
// properties of what it means to be a ShinySword. An actual
// ShinySword lying around is an Instance.
type Definition struct {
	id                  string
	Aliases             []string
	Categories          CategorySet
	Name                string
	ShortDescription    string // description of the object when being used: "a long, green stick" -> "The Beastly Fido picks up the long, green stick."
	DescriptionOnGround string // description of the object when lying on the ground: "A shiny sword is lying here."
	ZoneId              string
	WearLocation        slot.Location
}

func BuildDefinitionId(zoneId string, definition string) string {
	return zoneId + ":" + definition
}

func NewDefinition(definitionId string, name string, zoneId string, category Category, aliases []string, shortDescription string, descriptionOnGround string, wearLocation slot.Location) *Definition {
	d := &Definition{
		id:                  definitionId,
		Name:                name,
		ShortDescription:    shortDescription,
		DescriptionOnGround: descriptionOnGround,
		Categories:          make(map[Category]bool),
		Aliases:             aliases,
		ZoneId:              zoneId,
		WearLocation:        wearLocation,
	}
	d.AddCategory(category)
	return d
}

func (d *Definition) Id() string {
	return BuildDefinitionId(d.ZoneId, d.id)
}

// return only the 'id' part of the Definition
func (d *Definition) Identifier() string {
	return d.id
}

func (d *Definition) AddCategory(cat Category) {
	d.Categories.Add(cat)
}

func (d *Definition) CanEquipWeapon() bool {
	// for now, you can equip this if it is a weapon
	return d.Categories.Contains(WEAPON)
}
