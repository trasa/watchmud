package object

import "fmt"

type InstanceNotFoundError struct {
	Id string
}

func (e *InstanceNotFoundError) Error() string {
	return fmt.Sprintf("Object Instance not found: %s", e.Id)
}
