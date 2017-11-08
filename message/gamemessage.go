package message

import (
	"errors"
	"fmt"
	"reflect"
)

//go:generate protoc -I . message.proto --go_out=plugins=grpc:.
//go:generate type-registry -in ./message.pb.go -typenameprefix GameMessage_

// make a game message that contains this inner object
func NewGameMessage(payload interface{}) (resultPtr *GameMessage, err error) {
	// our result values
	result := GameMessage{}
	resultPtr = &result

	// Get the "Inner" field of the GameMessage struct
	// so we can populate it
	mValue := reflect.ValueOf(resultPtr)
	gamemessageInnerField := mValue.Elem().FieldByName("Inner")

	// wrapperType is something like GameMessage_DropResponse
	// look up the GameMessage_* type from innerTypes based on the type of payload
	wrapperType, ok := typeregistry[DecodeTypeName(payload)]
	if !ok {
		err = errors.New(fmt.Sprintf("message.NewGameMessage: Can't find wrapper type for %s", reflect.TypeOf(payload).String()))
		return
	}

	// create GameMessage_*** from the wrapperType
	wrapper := reflect.New(wrapperType)

	// GameMessage_[Response] has a field called "Response"
	// example: GameMessage_DropResponse.DropResponse
	// Find the field based on the name of the payload
	var innerValueType reflect.Type
	var payloadValue reflect.Value
	payloadType := reflect.TypeOf(payload)
	if payloadType.Kind() == reflect.Ptr {
		// payload is a ptr, so get the type of what payload points to
		innerValueType = reflect.TypeOf(payload).Elem()
		// also need payload == the value of payload's ptr
		payloadValue = reflect.ValueOf(payload)
	} else {
		// payload is not a ptr, just use the payload itself to determine type
		innerValueType = reflect.TypeOf(payload)
		// however the value does need to be *payload
		payloadValue = reflect.New(payloadType)
		payloadValue.Elem().Set(reflect.ValueOf(payload))
	}

	// "DropResponse"
	innerValueFieldName := innerValueType.Name()
	// GameMessage_DropResponse.DropResponse
	innerValueField := wrapper.Elem().FieldByName(innerValueFieldName)

	// set wrapper.GameMessage_payloadType = payload object
	// wrapper.GameMessage_payloadType is a *Payload
	innerValueField.Set(payloadValue)

	// set the innerField to our wrapper object
	gamemessageInnerField.Set(wrapper)

	return
}
