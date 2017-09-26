package message

import (
	"encoding/json"
	"errors"
	"github.com/mitchellh/mapstructure"
	"github.com/trasa/watchmud/direction"
	"log"
	"reflect"
	"strings"
)

const ENABLE_DECODE_LOGGING = false

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
	case "drop":
		request = DropRequest{
			Request: RequestBase{MessageType: "drop"},
			Target:  tokens[1],
		}

	case "exits", "exit", "ex":
		request = ExitsRequest{
			Request: RequestBase{MessageType: "exits"},
		}

	case "get":
		request = GetRequest{
			Request: RequestBase{MessageType: "get"},
			Targets: tokens[1:],
		}

	case "inv", "inventory":
		request = InventoryRequest{
			Request: RequestBase{MessageType: "inv"},
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
// this is deprecated in favor of line input, other than for login.
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

	//noinspection GoBoolExpressions
	if ENABLE_DECODE_LOGGING {
		log.Printf("translateToResponse: raw: %s", string(raw))
	}

	responseMap := rawMap["Response"].(map[string]interface{})
	//noinspection GoNameStartsWithPackageName
	messageType := responseMap["msg_type"].(string)

	switch messageType {
	case "drop":
		response = decodeResponse(&DropResponse{}, rawMap)

	case "enter_room":
		response = decodeResponse(&EnterRoomNotification{}, rawMap)

	case "error":
		response = decodeResponse(&ErrorResponse{}, rawMap)

	case "exits":
		response = decodeResponse(&ExitsResponse{}, rawMap)

	case "get":
		response = decodeResponse(&GetResponse{}, rawMap)

	case "inv":
		response = decodeResponse(&InventoryResponse{}, rawMap)

	case "leave_room":
		response = decodeResponse(&LeaveRoomNotification{}, rawMap)

	case "login_response":
		response = decodeResponse(&LoginResponse{}, rawMap)

	case "look":
		response = decodeResponse(&LookResponse{}, rawMap)

	case "move":
		response = decodeResponse(&MoveResponse{}, rawMap)

	case "say":
		response = decodeResponse(&SayResponse{}, rawMap)

	case "say_notification":
		response = decodeResponse(&SayNotification{}, rawMap)

	case "tell":
		response = decodeResponse(&TellResponse{}, rawMap)

	case "tell_notification":
		response = decodeResponse(&TellNotification{}, rawMap)

	case "tell_all":
		response = decodeResponse(&TellAllResponse{}, rawMap)

	case "tell_all_notification":
		response = decodeResponse(&TellAllNotification{}, rawMap)

	case "who":
		response = decodeResponse(&WhoResponse{}, rawMap)

	default:
		err = &UnknownMessageTypeError{MessageType: messageType}
		log.Printf("translator.TranslateToResponse: unknown message type: %v", err)
		return
	}

	if err != nil {
		log.Println("Failed to decode rawMap:", rawMap, err)
		return
	}
	fillResponseBase(response, responseMap)
	return
}

// take rawmap and use it to create a Response
// decode the map into the response structure
func decodeResponse(response interface{}, rawMap map[string]interface{}) Response {

	// (no kidding this constant is always false or true ... )
	//noinspection ALL
	if ENABLE_DECODE_LOGGING {
		log.Println("method call type ", reflect.TypeOf(response))
		log.Println("method call response kind", reflect.TypeOf(response).Kind())
	}
	mapstructure.Decode(rawMap, response)
	return response.(Response)
}

// set the ResponseBase members (they aren't set through mapstructure.Decode)
func fillResponseBase(response Response, responseMap map[string]interface{}) {
	//noinspection GoNameStartsWithPackageName
	messageType := responseMap["msg_type"].(string)
	reflect.ValueOf(response).Elem().FieldByName("Response").Set(reflect.ValueOf(NewSuccessfulResponse(messageType)))
	response.SetResultCode(responseMap["result_code"].(string))
	response.SetSuccessful(responseMap["success"].(bool))
	response.SetMessageType(responseMap["msg_type"].(string))
}

func TranslateToJson(obj interface{}) (result string, err error) {
	raw, err := json.Marshal(obj)
	result = string(raw)
	return
}
