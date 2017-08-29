package message

type Response interface {
	GetMessageType() string
	IsSuccessful() bool
	GetResultCode() string
}

type ResponseBase struct {
	MessageType string `json:"msg_type"`
	Successful  bool   `json:"success"`
	ResultCode  string `json:"result_code"`
}

func (r ResponseBase) GetMessageType() string {
	return r.MessageType
}

func (r ResponseBase) IsSuccessful() bool {
	return r.Successful
}

func (r ResponseBase) GetResultCode() string {
	return r.ResultCode
}

func NewSuccessfulResponse(msgType string) Response {
	return ResponseBase{
		MessageType: msgType,
		Successful:  true,
		ResultCode:  "OK",
	}
}
