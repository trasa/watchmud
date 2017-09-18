package message

type WhoRequest struct {
	Request
}

type WhoResponse struct {
	Response
	PlayerInfo []WhoPlayerInfo
}

type WhoPlayerInfo struct {
	PlayerName string
	ZoneName   string
	RoomName   string
}
