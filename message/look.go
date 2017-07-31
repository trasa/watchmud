package message

type LookRequest struct {
	Request
	ValueList []string `json:"value_list"`
}

type LookResponse struct {
	Response
	RoomName string `json:"room_name"`
	Value    string `json:"value"`
	// TODO MOBs, Players, Objects ....
}
