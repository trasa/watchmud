package message

import (
	"reflect"
	"strings"
)

// Turn payload type.String() (which 'does the right thing' if
// payload is a pointer type) into just the name of the type,
// minus package name, minus any "*" if needed.
func DecodeTypeName(payload interface{}) string {
	if payload == nil {
		return ""
	}
	s := strings.Replace(reflect.TypeOf(payload).String(), "*", "", -1)
	s = strings.Replace(s, "message.", "", -1)
	return s
}
