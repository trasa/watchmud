package client

import "github.com/trasa/watchmud/player"

// noinspection GoNameStartsWithPackageName
type Client interface {
	Player() *player.Player
	Send(innerMessage any) error
	Close()
}
