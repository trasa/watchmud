package client

import "github.com/trasa/watchmud/player"

//noinspection GoNameStartsWithPackageName
type Client interface {
	Send(message interface{})
	SetPlayer(player player.Player)
	GetPlayer() player.Player
	Close()
}
