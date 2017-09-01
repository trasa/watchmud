package message

import (
	"encoding/json"
	"errors"
	"github.com/mitchellh/mapstructure"
	"github.com/trasa/watchmud/direction"
	"log"
	"strings"
)

// Turn a request of format type "line" into a Request message
//
// have an input string like
// tell bob hi there
// turn into Request = message.TellRequest { "bob", "hi there" }
// and so on.
// note that not all commands can be parsed from line input
func TranslateLineToRequest(line string) (request Request, err error) {

	// first parse into tokens
	tokens := strings.Fields(line)
	if len(tokens) == 0 {
		request = NoOpRequest{
			Request: RequestBase{MessageType: "no_op"},
		}
		return
	}
	switch tokens[0] {
	case "exits", "ex":
		request = ExitsRequest{
			Request: RequestBase{MessageType: "exits"},
		}
	case "look", "l":
		request = LookRequest{
			Request:   RequestBase{MessageType: "look"},
			ValueList: tokens[1:],
		}
	case "n", "north", "s", "south", "e", "east", "w", "west", "u", "up", "d", "down":
		var d direction.Direction
		if d, err = direction.StringToDirection(tokens[0]); err == nil {
			request = MoveRequest{
				Request:   RequestBase{MessageType: "move"},
				Direction: d,
			}
		}
	case "'", "say":
		if len(tokens) >= 2 {
			request = SayRequest{
				Request: RequestBase{MessageType: "say"},
				Value:   strings.Join(tokens[1:], " "),
			}
		} else {
			err = errors.New("What do you want to say?")
		}
	case "tell", "t":
		if len(tokens) >= 3 {
			request = TellRequest{
				Request:            RequestBase{MessageType: "tell"},
				ReceiverPlayerName: tokens[1],
				Value:              strings.Join(tokens[2:], " "),
			}
		} else {
			// some sort of error about malformed tell request...
			err = errors.New("usage: tell [somebody] [something]")
		}
	case "tellall", "ta":
		if len(tokens) >= 2 {
			request = TellAllRequest{
				Request: RequestBase{MessageType: "tell_all"},
				Value:   strings.Join(tokens[1:], " "),
			}
		} else {
			err = errors.New("usage: tellall [something]")
		}
	case "who":
		request = WhoRequest{
			Request: RequestBase{MessageType: "who"},
		}
	default:
		err = errors.New("Unknown request: " + tokens[0])
	}
	return
}

// Turn a request of format type "request" into a Request object
func TranslateToRequest(body map[string]interface{}) (request Request, err error) {
	err = nil
	msgType := body["Request"].(map[string]interface{})["msg_type"].(string)
	requestBase := RequestBase{MessageType: msgType}

	switch msgType {
	case "login":
		//{
		// 	"Request":{
		// 			"msg_type":"login"
		// 		},
		// 	"player_name":"somedood",
		// 	"password":"NotImplemented"
		// }
		var lr LoginRequest
		err = mapstructure.Decode(body, &lr)
		lr.Request = requestBase
		request = lr

	case "tell":
		var tr TellRequest
		err = mapstructure.Decode(body, &tr)
		tr.Request = requestBase
		request = tr

	case "tell_all":
		var tar TellAllRequest
		err = mapstructure.Decode(body, &tar)
		tar.Request = requestBase
		request = tar

	default:
		err = &UnknownMessageTypeError{MessageType: msgType}
	}
	return
}

func TranslateToResponse(raw []byte) (response Response, err error) {
	var rawMap map[string]interface{}
	if err = json.Unmarshal(raw, &rawMap); err != nil {
		log.Println("Unmarshal error: ", err)
		return
	}
	responseMap := rawMap["Response"].(map[string]interface{})
	switch responseMap["msg_type"].(string) {
	case "login_response":
		// allocate a LoginResponse and also a loginResponse.Response, or
		// things will fail later on (must do this for all msg_types)
		loginResp := &LoginResponse{
			Response: NewSuccessfulResponse("login_response"),
		}
		mapstructure.Decode(rawMap, &loginResp)
		response = loginResp

	default:
		err = &UnknownMessageTypeError{MessageType: rawMap["msg_type"].(string)}
		log.Println("unknown message type: ", err)
		return
	}

	// set the ResponseBase (doesn't work through mapstructure.Decode)
	// note that response.Response must have been allocated above in the
	// switch, or this code will fail with a nil pointer.
	response.SetResultCode(responseMap["result_code"].(string))
	response.SetSuccessful(responseMap["success"].(bool))
	response.SetMessageType(responseMap["msg_type"].(string))

	return
}
