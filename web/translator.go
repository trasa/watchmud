package web

import (
	"errors"
	"github.com/trasa/watchmud/direction"
	"github.com/trasa/watchmud/message"
	"strings"
)

// have an input string like
// tell bob hi there
// turn into Request = message.TellRequest { "bob", "hi there" }
// and so on.
// note that not all commands can be parsed from line input
func translateLineToRequest(line string) (request message.Request, err error) {

	// first parse into tokens
	tokens := strings.Fields(line)
	if len(tokens) == 0 {
		request = message.NoOpRequest{
			Request: message.RequestBase{MessageType: "no_op"},
		}
		return
	}
	switch tokens[0] {
	case "look", "l":
		request = message.LookRequest{
			Request:   message.RequestBase{MessageType: "look"},
			ValueList: tokens[1:],
		}
	case "n", "north", "s", "south", "e", "east", "w", "west", "u", "up", "d", "down":
		var d direction.Direction
		if d, err = direction.StringToDirection(tokens[0]); err == nil {
			request = message.GoRequest{
				Request:   message.RequestBase{MessageType: "go"},
				Direction: d,
			}
		}
	case "tell", "t":
		if len(tokens) >= 3 {
			request = message.TellRequest{
				Request:            message.RequestBase{MessageType: "tell"},
				ReceiverPlayerName: tokens[1],
				Value:              strings.Join(tokens[2:], " "),
			}
		} else {
			// some sort of error about malformed tell request...
			err = errors.New("usage: tell [somebody] [something]")
		}
	case "tellall", "ta":
		if len(tokens) >= 2 {
			request = message.TellAllRequest{
				Request: message.RequestBase{MessageType: "tell_all"},
				Value:   strings.Join(tokens[1:], " "),
			}
		} else {
			err = errors.New("usage: tellall [something]")
		}
	default:
		err = errors.New("Unknown command: " + tokens[0])
	}
	return
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
