package demo

import (
	"context"
	"github.com/jzyong/GameServer4g/game-message/message"
	grpc "google.golang.org/grpc"
	"log"
	"testing"
	"time"
)

func TestClientDial(t *testing.T) {
	conn, err := grpc.Dial("localhost:1234", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer conn.Close()

	client := message.NewServerServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	result, err := client.ServerRegister(ctx, &message.ServerInfo{
		Id:    1,
		Ip:    "1",
		Type:  0,
		State: 0,
	})
	if err != nil {
		log.Fatalf("%v", err)
	}
	log.Printf("%s", result)
}
