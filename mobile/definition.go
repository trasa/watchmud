package mobile

// Defines what it means to be a mob.
type Definition struct {
	DefinitionId      string
	Aliases           []string
	Name              string
	ShortDescription  string
	DescriptionInRoom string // description when in a room "A giant lizard is here."
	CanWander         bool   // walks from room to room
}

func NewDefinition(definitionId string, name string, aliases []string, shortDescription, descriptionInRoom string, canWander bool) *Definition {
	d := &Definition{
		DefinitionId:      definitionId,
		Name:              name,
		Aliases:           aliases,
		ShortDescription:  shortDescription,
		DescriptionInRoom: descriptionInRoom,
		CanWander:         canWander,
	}
	return d
}
