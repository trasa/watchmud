package web

import "fmt"

type UnknownMessageTypeError struct {
	MessageType string
}

func (e *UnknownMessageTypeError) Error() string {
	return fmt.Sprintf("Unknown MessageType: %s", e.MessageType)
}
