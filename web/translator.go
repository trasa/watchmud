package web

import (
	"errors"
	"github.com/trasa/watchmud/message"
)

func translateLineToRequest(line string) (request message.Request, err error) {
	// have an input string like
	// tell bob hi there
	// turn into Request = message.TellRequest { "bob", "hi there" }
	// and so on.
	// difficulty: don't repeat yourself with the translateToRequest below that turns a map
	// into an object ... so should this return a map or jump directly to Request?
	return nil, errors.New("not implemented")
}
func translateToRequest(body map[string]string) (request message.Request, err error) {
	err = nil
	switch body["msg_type"] {
	case "login":
		request = message.LoginRequest{
			Request:    message.RequestBase{MessageType: body["msg_type"]},
			PlayerName: body["player_name"],
			Password:   body["password"],
		}
	case "tell":
		request = message.TellRequest{
			Request:            message.RequestBase{MessageType: body["msg_type"]},
			ReceiverPlayerName: body["receiver"],
			Value:              body["value"],
		}
	case "tell_all":
		request = message.TellAllRequest{
			Request: message.RequestBase{MessageType: body["msg_type"]},
			Value:   body["value"],
		}
	default:
		err = &UnknownMessageTypeError{MessageType: body["msg_type"]}
	}
	return
}
