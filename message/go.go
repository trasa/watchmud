package message

import (
	"github.com/trasa/watchmud/direction"
)

type GoRequest struct {
	Request
	Direction direction.Direction `json:"direction"`
}

type GoResponse struct {
	Response
}
