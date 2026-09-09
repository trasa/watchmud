package object

import "fmt"

type Id struct {
	DefinitionId string
	ZoneId       string
}

func NewObjectId(id string, zoneId string) Id {
	return Id{
		DefinitionId: id,
		ZoneId:       zoneId,
	}
}

func (id *Id) String() string {
	return fmt.Sprintf("%s:%s", id.ZoneId, id.DefinitionId)
}
