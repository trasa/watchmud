package spaces

import "fmt"

/*
// The known set of zone commands
// Used for storing  the command information
type CommandType int

const (
	UNKNOWN CommandType = iota
	CREATE_OBJECT
	CREATE_MOBILE
	GIVE_OBJECT_TO_MOBILE
	EQUIP_OBJECT_ON_MOBILE
	PUT_OBJECT_IN_OBJECT
	DOOR_STATE
	REMOVE_OBJECT_FROM_ROOM
)
*/

// Supertype for all Zone Commands
type ZoneCommand interface {
}

// CreateObject instructs the world to create an object in the world.
type CreateObject struct {
	ObjectDefinitionId string // what type of object
	ZoneId             string // which zone is the definition in, leave empty for "this zone"
	RoomId             string // where does the object go?
	InstanceMax        int    // how many are allowed to be lying around the zone?
}

func (cmd CreateObject) String() string {
	return fmt.Sprintf("Create Object '%s-%s' in Room '%s', Max of %d",
		cmd.ObjectDefinitionId,
		cmd.ZoneId,
		cmd.RoomId,
		cmd.InstanceMax,
	)
}

// Command: Create a Mobile
type CreateMobile struct {
	MobileDefinitionId string // what type of mobile
	ZoneId             string // where the mobile is defined, or empty for "this zone"
	RoomId             string // where does the mobile go?
	InstanceMax        int    // how many are allowed to be walking around the zone?
	// TODO give equipment
	// TODO give objects
}

func (cmd CreateMobile) String() string {
	return fmt.Sprintf("Create Mobile '%s-%s' in Room '%s', Max of %d",
		cmd.MobileDefinitionId,
		cmd.ZoneId,
		cmd.RoomId,
		cmd.InstanceMax,
	)
}

// TODO: other types
