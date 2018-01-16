package object

import "fmt"

type InstanceNotWeaponError struct {
	Id string
}

func (e *InstanceNotWeaponError) Error() string {
	return fmt.Sprintf("Object Instance not a weapon: %s", e.Id)
}
