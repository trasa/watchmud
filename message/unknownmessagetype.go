package message

import "fmt"

// TODO do we still need this?
type UnknownMessageTypeError struct {
	MessageType string
}

func (e *UnknownMessageTypeError) Error() string {
	return fmt.Sprintf("Unknown MessageType: %s", e.MessageType)
}
