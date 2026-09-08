package gameserver

type Instance interface {
	Receive(handlerParam *HandlerParameter)
	Logout(c Conn, cause string)
}
