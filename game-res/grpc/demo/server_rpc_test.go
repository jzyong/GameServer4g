package demo

import (
	"context"
	"github.com/jzyong/GameServer4g/game-message/message"
	"google.golang.org/grpc"
	"log"
	"net"
	"testing"
)

func TestStartServer(t *testing.T) {

	grpcServer := grpc.NewServer()
	message.RegisterServerServiceServer(grpcServer, new(ServerServiceImpl))
	listen, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal(err)
	}
	grpcServer.Serve(listen)
}

type ServerServiceImpl struct {
}

func (server *ServerServiceImpl) ServerRegister(ctx context.Context, in *message.ServerInfo) (*message.ServerInfo, error) {
	log.Printf("%v", in)
	return in, nil
}

func (server *ServerServiceImpl) ServerUpdate(ctx context.Context, in *message.ServerInfo) (*message.ServerInfo, error) {
	log.Printf("%v", in)
	return in, nil
}
