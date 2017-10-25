package object

// Definition of what it means to be an "object"
// the "platonic form" of an object, if you will.
//
// There is a definition of the type "ShinySword" which defines the
// properties of what it means to be a ShinySword. An actual
// ShinySword lying around is an Instance.
type Definition struct {
	Id                  string
	Aliases             []string
	Categories          CategorySet
	Name                string
	ShortDescription    string // description of the object when being used: "a long, green stick" -> "The Beastly Fido picks up the long, green stick."
	DescriptionOnGround string // description of the object when lying on the ground: "A shiny sword is lying here."
	ZoneId              string
}

func NewDefinition(definitionId string, name string, zoneId string, category Category, aliases []string, shortDescription string, descriptionOnGround string) *Definition {
	d := &Definition{
		Id:                  definitionId,
		Name:                name,
		ShortDescription:    shortDescription,
		DescriptionOnGround: descriptionOnGround,
		Categories:          make(map[Category]bool),
		Aliases:             aliases,
		ZoneId:              zoneId,
	}
	d.AddCategory(category)
	return d
}

func (d *Definition) AddCategory(cat Category) {
	d.Categories.Add(cat)
}
