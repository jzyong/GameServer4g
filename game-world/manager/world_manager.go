package manager

import (
	"encoding/json"
	"fmt"
	"github.com/go-zookeeper/zk"
	game_common "github.com/jzyong/GameServer4g/game-common"
	"github.com/jzyong/GameServer4g/game-message/message"
	"github.com/jzyong/GameServer4g/game-world/config"
	"github.com/jzyong/GameServer4g/game-world/rpc"
	"github.com/jzyong/golib/log"
	"github.com/jzyong/golib/util"
	"google.golang.org/grpc"
	"net"
)

//世界
type WorldManager struct {
	util.DefaultModule
	ZKConnect  *zk.Conn //zookeeper连接
	GrpcServer *grpc.Server
}

var worldManager = &WorldManager{}

func GetWorldManager() *WorldManager {
	return worldManager
}

//
func (m *WorldManager) Init() error {
	log.Info("WorldManager:init")
	//初始化id
	util.UUID = util.NewSnowflake(int16(config.ApplicationConfigInstance.Id))

	// zookeeper 初始化
	//推送配置
	config := config.ApplicationConfigInstance
	m.ZKConnect = util.ZKCreateConnect(config.ZookeeperUrls)
	configBytes, _ := json.Marshal(config)
	util.ZKUpdate(m.ZKConnect, fmt.Sprintf(game_common.WorldConfig, config.Profile, config.Id), string(configBytes))

	// 启动grpc服务
	go m.startGrpcService()

	//注册服务
	util.ZKAdd(m.ZKConnect, fmt.Sprintf(game_common.WorldRpcService, config.Profile, config.Id), config.RpcUrl, zk.FlagEphemeral)

	log.Info("WorldManager:inited")
	return nil
}

//启动grpc
func (m *WorldManager) startGrpcService() {
	m.GrpcServer = grpc.NewServer()
	message.RegisterPlayerWorldServiceServer(m.GrpcServer, new(rpc.PlayerServiceImpl))
	listen, err := net.Listen("tcp", config.ApplicationConfigInstance.RpcUrl)
	if err != nil {
		log.Fatal("%v", err)
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
