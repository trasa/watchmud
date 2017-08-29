package message

type LookRequest struct {
	Request
	ValueList []string `json:"value_list"`
}

type LookResponse struct {
	Response
	RoomDescription
}

type LookNotification struct {
	Response
	RoomDescription
}

type RoomDescription struct {
	Name        string   `json:"room_name"`
	Description string   `json:"description"`
	Exits       string   `json:"exits"`
	Players     []string `json:"players"`
	// TODO MOBs, Objects ....
}
