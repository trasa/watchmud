package message


type LoginRequest struct {
	Request
	PlayerName string `json:"player_name"`
	Password string `json:"password"`
}


type LoginResponse struct {
	Response
	Player PlayerData `json:"player"`
}