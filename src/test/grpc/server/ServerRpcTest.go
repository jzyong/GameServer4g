package main

import (
	"context"
	"github.com/jzyong/go-mmo-server/src/message"
	"google.golang.org/grpc"
	"log"
	"net"
)

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

func main() {

	grpcServer := grpc.NewServer()
	message.RegisterServerServiceServer(grpcServer, new(ServerServiceImpl))
	listen, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal(err)
	}
	grpcServer.Serve(listen)
}
