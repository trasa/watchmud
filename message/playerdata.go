package message

type PlayerData struct {
	Name string
}

func NewPlayerData(name string) PlayerData {
	return PlayerData{
		Name: name,
	}
}
