package web

import (
	"errors"
	"github.com/mitchellh/mapstructure"
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
	case "exits", "ex":
		request = message.ExitsRequest{
			Request: message.RequestBase{MessageType: "exits"},
		}
	case "look", "l":
		request = message.LookRequest{
			Request:   message.RequestBase{MessageType: "look"},
			ValueList: tokens[1:],
		}
	case "n", "north", "s", "south", "e", "east", "w", "west", "u", "up", "d", "down":
		var d direction.Direction
		if d, err = direction.StringToDirection(tokens[0]); err == nil {
			request = message.MoveRequest{
				Request:   message.RequestBase{MessageType: "move"},
				Direction: d,
			}
		}
	case "'", "say":
		if len(tokens) >= 2 {
			request = message.SayRequest{
				Request: message.RequestBase{MessageType: "say"},
				Value:   strings.Join(tokens[1:], " "),
			}
		} else {
			err = errors.New("What do you want to say?")
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
	case "who":
		request = message.WhoRequest{
			Request: message.RequestBase{MessageType: "who"},
		}
	default:
		err = errors.New("Unknown request: " + tokens[0])
	}
	return
}

func translateToRequest(body map[string]interface{}) (request message.Request, err error) {
	err = nil
	msgType := body["Request"].(map[string]interface{})["msg_type"].(string)
	switch msgType {
	case "login":
		//{
		// 	"Request":{
		// 			"msg_type":"login"
		// 		},
		// 	"player_name":"somedood",
		// 	"password":"NotImplemented"
		// }
		var lr message.LoginRequest
		err = mapstructure.Decode(body, &lr)
		lr.Request = message.RequestBase{MessageType: msgType}
		request = lr

	case "tell":
		request = message.TellRequest{
			Request:            message.RequestBase{MessageType: msgType},
			ReceiverPlayerName: body["receiver"].(string),
			Value:              body["value"].(string),
		}
	case "tell_all":
		request = message.TellAllRequest{
			Request: message.RequestBase{MessageType: msgType},
			Value:   body["value"].(string),
		}
	default:
		err = &UnknownMessageTypeError{MessageType: msgType}
	}
	return
}
