package message

type LookRequest struct {
	Request
	ValueList []string `json:"value_list"`
}

type LookResponse struct {
	Response
	RoomDescription RoomDescription
}

type LookNotification struct {
	Response
	RoomDescription RoomDescription
}

type RoomDescription struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Exits       string   `json:"exits"`
	Players     []string `json:"players"`
	// TODO MOBs, Objects ....
}
