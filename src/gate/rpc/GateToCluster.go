package rpc

import (
	"context"
	"gate/config"
	"google.golang.org/grpc"
	"log"
	"message"
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
		Id:                   config.GateConfig.Id,
		BelongID:             config.GateConfig.Id,
		Ip:                   "localhost",
		Type:                 2,
		Port:                 80,
		State:                1,
		Version:              "",
		Content:              "gate",
		Online:               0,
		MaxUserCount:         2000,
		HttpPort:             80,
		OpenTime:             "",
		MaintainTime:         "",
		Name:                 "client",
		Wwwip:                "",
		XXX_NoUnkeyedLiteral: struct{}{},
		XXX_unrecognized:     nil,
		XXX_sizecache:        0,
	})
	if err != nil {
		log.Fatal("%v", err)
	}
	log.Printf("%s", result)
}
