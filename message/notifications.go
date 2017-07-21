package message

type Notification struct {
	MessageType string `json:"msg_type"`
}

type TellNotification struct {
	Notification
	From  string `json:"from"`
	Value string `json:"value"`
}
