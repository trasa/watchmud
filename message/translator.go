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

// map messageTypes to pointers to base Response objects
// these pointers will be used to create new Responses
// (don't use these directly!)
var responseRegistry = map[string]Response{
	"enter_room": &EnterRoomNotification{},
	"error": &ErrorResponse{},
	"exits": &ExitsResponse{},
	"login_response": &LoginResponse{},
	"look": &LookResponse{},
	"move": &MoveResponse{},
	"say": &SayResponse{},
}

func TranslateToResponse(raw []byte) (response Response, err error) {
	var rawMap map[string]interface{}
	if err = json.Unmarshal(raw, &rawMap); err != nil {
		log.Println("Unmarshal error: ", err)
		return
	}
	responseMap := rawMap["Response"].(map[string]interface{})
	//noinspection GoNameStartsWithPackageName
	messageType := responseMap["msg_type"].(string)


	//responsePrototype := responseRegistry[messageType]
	// create a new response to decode into
	// TODO
	// then fill the obj
	//response = decodeResponse(response)

	switch messageType {
	case "enter_room":
		response = decodeResponse(&EnterRoomNotification{}, rawMap)

	case "error":
		response = decodeResponse(&ErrorResponse{}, rawMap)

	case "exits":
		response = decodeResponse(&ExitsResponse{}, rawMap)

	case "login_response":
		response = decodeResponse(&LoginResponse{}, rawMap)

	case "look":
		response = decodeResponse(&LookResponse{}, rawMap)

	case "move":
		response = decodeResponse(&MoveResponse{}, rawMap)

	case "say":
		// can we get this to work using the existing type from the map?
		responsePrototypePtr := responseRegistry["say"]
		log.Printf("responsePrototypePtr %s is at addr %p", reflect.TypeOf(responsePrototypePtr), responsePrototypePtr)
		// allocate a new response based off of prototype
		// this allocates a **Response
		responseObjPtr := reflect.New(reflect.TypeOf(&responsePrototypePtr)).Interface()
		log.Printf("responseObjPtr type %s, string %s", reflect.TypeOf(responseObjPtr), responseObjPtr)

		// now need to allocate a Response
		responseObj := reflect.New(reflect.TypeOf(responseObjPtr))
		//responseObj = reflect.New(reflect.TypeOf(responseObjPtr.Elem().Interface().(Response)))

		log.Printf("responseObj type %s at addr %p", reflect.TypeOf(responseObj), responseObj)
		log.Printf("&responseObj type %s", reflect.TypeOf(&responseObj))
		log.Println("what im trying to get:", reflect.TypeOf(&SayResponse{}))
		log.Printf("does this equal? %s", reflect.TypeOf(&SayResponse{}) == reflect.TypeOf(responseObj))
		if reflect.TypeOf(&SayResponse{}) != reflect.TypeOf(responseObj) {
			panic("types are wrong")
		}

		log.Println("responseObj kind", reflect.TypeOf(responseObj).Kind())
		response = decodeResponse(responseObj, rawMap)
	default:
		err = &UnknownMessageTypeError{MessageType: messageType}
		log.Println("unknown message type: ", err)
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
func decodeResponse(response interface{}, rawMap map[string]interface{}) Response {
	// decode the map into the response structure
	log.Println("method call type ", reflect.TypeOf(response))
	log.Println("method call response kind", reflect.TypeOf(response).Kind())
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
