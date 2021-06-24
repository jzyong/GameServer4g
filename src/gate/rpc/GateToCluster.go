package rpc

import (
	"context"
	"github.com/jzyong/go-mmo-server/src/gate/config"
	"github.com/jzyong/go-mmo-server/src/message"
	"google.golang.org/grpc"
	"log"
	"time"
)

var GateToClusterClient *GateToCluster
var connect grpc.ClientConnInterface
var client message.ServerServiceClient

//连接cluster
type GateToCluster struct {
	connect grpc.ClientConnInterface
}

func (g *GateToCluster) Start(rpcUrl string) {
	conn, err := grpc.Dial(rpcUrl, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("%v", err)
	}
	connect = conn
	client = message.NewServerServiceClient(connect)
}

func (g *GateToCluster) Stop() {
	//connect.Close()
}

/**
注册到cluster服务器
*/
func RegisterToCluster() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	result, err := client.ServerRegister(ctx, &message.ServerInfo{
		Id:           config.GateConfigInstance.Id,
		BelongID:     config.GateConfigInstance.Id,
		Ip:           "localhost",
		Type:         2,
		Port:         80,
		State:        1,
		Version:      "",
		Content:      "gate",
		Online:       0,
		MaxUserCount: 2000,
		HttpPort:     80,
		OpenTime:     "",
		MaintainTime: "",
		Name:         "client",
		Wwwip:        "",
	})
	if err != nil {
		log.Printf("%v", err)
	}
	log.Printf("%s", result)
}
