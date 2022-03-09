package manager

import (
	"encoding/json"
	"fmt"
	"github.com/go-zookeeper/zk"
	"github.com/jzyong/go-mmo-server/src/core/log"
	"github.com/jzyong/go-mmo-server/src/core/util"
	"github.com/jzyong/go-mmo-server/src/message"
	"github.com/jzyong/go-mmo-server/src/world/config"
	"github.com/jzyong/go-mmo-server/src/world/rpc"
	"google.golang.org/grpc"
	"net"
)

//世界
type WorldManager struct {
	util.DefaultModule
	ZKConnect  *zk.Conn //zookeeper连接
	GrpcServer *grpc.Server
}

func NewWorldManager() *WorldManager {
	return &WorldManager{}
}

//
func (m *WorldManager) Init() error {
	log.Info("WorldManager:init")
	//初始化id
	util.UUID = util.NewSnowflake(int16(config.WorldConfigInstance.Id))

	// zookeeper 初始化
	//推送配置
	config := config.WorldConfigInstance
	m.ZKConnect = util.ZKCreateConnect(config.ZookeeperUrls)
	configBytes, _ := json.Marshal(config)
	util.ZKUpdate(m.ZKConnect, fmt.Sprintf(util.WorldConfig, config.Profile, config.Id), string(configBytes))

	// 启动grpc服务
	go m.startGrpcService()

	//注册服务
	util.ZKAdd(m.ZKConnect, fmt.Sprintf(util.WorldRpcService, config.Profile, config.Id), config.RpcUrl, zk.FlagEphemeral)

	log.Info("WorldManager:inited")
	return nil
}

//启动grpc
func (m *WorldManager) startGrpcService() {
	m.GrpcServer = grpc.NewServer()
	message.RegisterPlayerWorldServiceServer(m.GrpcServer, new(rpc.PlayerServiceImpl))
	listen, err := net.Listen("tcp", config.WorldConfigInstance.RpcUrl)
	if err != nil {
		log.Fatal(err)
	}
	m.GrpcServer.Serve(listen)
}

func (m *WorldManager) Stop() {
	if m.ZKConnect != nil {
		m.ZKConnect.Close()
	}
	if m.GrpcServer != nil {
		m.GrpcServer.Stop()
	}
}
