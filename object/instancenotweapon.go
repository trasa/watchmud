package object

import "fmt"

// TODO replace with error
type InstanceNotWeaponError struct {
	Id string
}

func (e *InstanceNotWeaponError) Error() string {
	return fmt.Sprintf("Object Instance not a weapon: %s", e.Id)
}
