package server

type Client interface {
	Send(message interface{}) // todo return err
	SetPlayer(player *Player)
	GetPlayer() *Player
	Close()
}

