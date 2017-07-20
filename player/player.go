package player

type Player interface {
	Send(message interface{}) // todo return err
	GetName() string
}
