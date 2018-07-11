package rpc

import (
	"fmt"
	"github.com/trasa/watchmud-message"
	"github.com/trasa/watchmud/gameserver"
	"google.golang.org/grpc"
	"log"
	"net"
)

type server struct {
	gameServerInstance gameserver.Instance
}

// Begin listening on address and port
func (s *server) Run(port int) {
	// TODO get from configuration
	log.Printf("gRPC listening on port %d", port)
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	message.RegisterMudCommServer(grpcServer, s)
	grpcServer.Serve(lis)
}

// While this method is open, we have a connection from the client to the server.
// When the client disconnects (or this returns), the connection is closed.
func (s *server) SendReceive(stream message.MudComm_SendReceiveServer) error {
	c := newClient(stream, s.gameServerInstance)
	go c.writePump()
	c.readPump()

	return nil
}

func NewServer(gs gameserver.Instance) *server {
	s := server{
		gameServerInstance: gs,
	}
	return &s
}
