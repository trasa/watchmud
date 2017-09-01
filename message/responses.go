package message

type Response interface {
	GetMessageType() string
	SetMessageType(string)

	IsSuccessful() bool
	SetSuccessful(bool)

	GetResultCode() string
	SetResultCode(string)
}

type ResponseBase struct {
	MessageType string `json:"msg_type"`
	Successful  bool   `json:"success"`
	ResultCode  string `json:"result_code"`
}

func (r *ResponseBase) SetMessageType(s string) {
	r.MessageType = s
}

func (r *ResponseBase) GetMessageType() string {
	return r.MessageType
}

func (r *ResponseBase) SetSuccessful(b bool) {
	r.Successful = b
}

func (r *ResponseBase) IsSuccessful() bool {
	return r.Successful
}

func (r *ResponseBase) SetResultCode(s string) {
	r.ResultCode = s
}

func (r *ResponseBase) GetResultCode() string {
	return r.ResultCode
}

func NewSuccessfulResponse(msgType string) Response {
	return &ResponseBase{
		MessageType: msgType,
		Successful:  true,
		ResultCode:  "OK",
	}
}

func NewUnsuccessfulResponse(msgType string, resultCode string) Response {
	return &ResponseBase{
		MessageType: msgType,
		Successful:  false,
		ResultCode:  resultCode,
	}
}
