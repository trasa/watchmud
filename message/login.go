package message

type LoginRequest struct {
	Request    `mapstructure:"-"`
	PlayerName string
	Password   string
}

type LoginResponse struct {
	ResponseBase
	Player PlayerData
}
