package message

type LookRequest struct {
	Request
	ValueList []string `json:"value_list"`
}

type LookResponse struct {
	Response
	RoomName    string `json:"room_name"`
	Description string `json:"description"`
	Exits       string `json:"exits"`
	// TODO MOBs, Players, Objects ....
}

type LookNotification struct {
	Notification
	RoomName    string `json:"room_name"`
	Description string `json:"description"`
	Exits       string `json:"exits"`
	// TODO objects and so on ...
}
