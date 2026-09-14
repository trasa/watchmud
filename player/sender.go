package player

// Sender is anything that can deliver a message to this player's connection
type Sender interface {
	Send(msg any)
}
