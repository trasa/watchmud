package server

var GameServerInstance *GameServer

type GameServer struct {
	incomingMessageBuffer chan *IncomingMessage
	World                 *World
}

func newGameServer() *GameServer {
	return &GameServer{
		incomingMessageBuffer: make(chan *IncomingMessage),
		World: NewWorld(),
	}
}

// Initialize the game server
func Init() {
	GameServerInstance = newGameServer()
}

func (server *GameServer) Run() {

	StartAllClientDispatcher()

	// this is the loop that handles incoming requests
	// needs to be organized around TICKs
	for {
		select {
		// TODO add in tick time
		case message := <-server.incomingMessageBuffer:
			server.World.handleIncomingMessage(message)
		}
	}
}
