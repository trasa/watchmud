package message

type Request interface {
	GetMessageType() string
}

type RequestBase struct {
	MessageType string `json:"msg_type"`
}

func (r RequestBase) GetMessageType() string {
	return r.MessageType
}
