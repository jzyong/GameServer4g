package main

import (
	"context"
	"github.com/jzyong/go-mmo-server/src/message"
	grpc "google.golang.org/grpc"
	"log"
	"time"
)

func main() {
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
		log.Fatal("%v", err)
	}
	log.Printf("%s", result)
}
