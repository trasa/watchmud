package message

type WhoRequest struct {
	Request
}

type WhoResponse struct {
	Response
	PlayerInfo []WhoPlayerInfo `json:"player_info"`
}

type WhoPlayerInfo struct {
	PlayerName string `json:"player_name"`
	ZoneName   string `json:"zone_name"`
	RoomName   string `json:"room_name"`
}
