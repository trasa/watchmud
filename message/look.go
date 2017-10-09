package message

type LookRequest struct {
	Request
	ValueList []string
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
	Name        string
	Description string
	Exits       string
	Players     []string
	Objects     []string
	Mobs        []string
}
