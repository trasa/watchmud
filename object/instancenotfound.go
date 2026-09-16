package object

import "fmt"

// TODO replace with error

type InstanceNotFoundError struct {
	Id string
}

func (e *InstanceNotFoundError) Error() string {
	return fmt.Sprintf("Object Instance not found: %s", e.Id)
}
