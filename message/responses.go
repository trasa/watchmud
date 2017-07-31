package message

type Response struct {
	MessageType string `json:"msg_type"`
	Successful  bool   `json:"success"`
	ResultCode  string `json:"result_code"`
}

func NewSuccessfulResponse(msgType string) Response {
	return Response{
		MessageType: msgType,
		Successful:  true,
		ResultCode:  "OK",
	}
}
