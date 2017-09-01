package message

import (
	"encoding/json"
	"log"
)

type Request interface {
	GetMessageType() string
}

type RequestBase struct {
	MessageType string `json:"msg_type" mapstructure:"-"`
}

func (r RequestBase) GetMessageType() string {
	return r.MessageType
}

// Raps up a Format identifier and a 'value' which is
// serialized by whatever 'Format' specifies.
// Where Value might be a serialized String, or a
// Map[string]interface{}
type RequestEnvelope struct {
	Format string      `json:"format"`
	Value  interface{} `json:"value"`
}

//
//func (r RequestEnv) MarshalJSON() error {
//	// TODO
//	return nil
//}
//

// TODO dead code?
func UnmarshalJSON(b []byte) error {
	// TODO
	log.Println("Unmarshal RequestEnv JSON:", string(b[:]))
	var raw interface{}
	err := json.Unmarshal(b, &raw)
	if err != nil {
		return err
	}

	return nil
}
