package message

type LoginRequest struct {
	Request    `mapstructure:"-"`
	PlayerName string
	Password   string
}

type LoginResponse struct {
	Response
	Player PlayerData
}
