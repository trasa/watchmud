package client

import "github.com/trasa/watchmud/player"

//noinspection GoNameStartsWithPackageName
type Client interface {
	Send(message interface{}) // todo return err
	SetPlayer(player player.Player)
	GetPlayer() player.Player
	Close()
}
