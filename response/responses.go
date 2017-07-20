package response

//noinspection GoNameStartsWithPackageName
type Response struct {
	MessageType string `json:"msg_type"`
	Successful  bool   `json:"success"`
	ResultCode  string `json:"result_code"`
}

type Notification struct {
	MessageType string `json:"msg_type"`
}

type LoginResponse struct {
	Response
	Player PlayerData `json:"player"`
}

type TellNotification struct {
	Notification
	From  string `json:"from"`
	Value string `json:"value"`
}

type TellAllNotification struct {
	Notification
	Value  string `json:"value"`
	Sender string `json:"sender"`
}
