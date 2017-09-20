package object

// Definition of what it means to be an "object"
// the "platonic form" of an object, if you will.
//
// There is a definition of the type "ShinySword" which defines the
// properties of what it means to be a ShinySword. An actual
// ShinySword lying around is an Instance.
type Definition struct {
	TypeId              string
	Aliases             []string
	Categories          map[Category]bool
	Name                string
	ShortDescription    string // description of the object when being used: "a long, green stick" -> "The Beastly Fido picks up the long, green stick."
	DescriptionOnGround string // description of the object when lying on the ground: "A shiny sword is lying here."

	// if true, multiple instances of this type can be "stacked" on top of each other
	// and don't need to be	individually accessed. Think "stack of papers", vs.
	// two shiny swords, one with a curse and one without...
	Stackable bool
}

func NewDefinition(typeId string, name string, category Category, stackable bool, aliases []string, shortDescription string, descriptionOnGround string) *Definition {
	d := &Definition{
		TypeId:              typeId,
		Name:                name,
		ShortDescription:    shortDescription,
		DescriptionOnGround: descriptionOnGround,
		Categories:          make(map[Category]bool),
		Aliases:             aliases,
		Stackable:           stackable,
	}
	d.AddCategory(category)
	return d
}

func (d *Definition) AddCategory(cat Category) {
	d.Categories[cat] = true
}
