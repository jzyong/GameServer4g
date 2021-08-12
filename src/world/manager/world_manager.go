package manager

import (
	"encoding/json"
	"fmt"
	"github.com/go-zookeeper/zk"
	"github.com/jzyong/go-mmo-server/src/core/log"
	"github.com/jzyong/go-mmo-server/src/core/util"
	"github.com/jzyong/go-mmo-server/src/world/config"
)

//世界
type WorldManager struct {
	util.DefaultModule
	ZKConnect *zk.Conn //zookeeper连接
}

func NewWorldManager() *WorldManager {
	return &WorldManager{}
}

//
func (this *WorldManager) Init() error {
	log.Info("WorldManager:init")
	//初始化id
	util.UUID = util.NewSnowflake(int16(config.WorldConfigInstance.Id))

	// zookeeper 初始化
	//推送配置
	config := config.WorldConfigInstance
	this.ZKConnect = util.ZKCreateConnect(config.ZookeeperUrls)
	configBytes, _ := json.Marshal(config)
	util.ZKUpdate(this.ZKConnect, fmt.Sprintf(util.WorldConfig, config.Profile, config.Id), string(configBytes))

	//TODO 启动grpc服务

	//注册服务
	util.ZKAdd(this.ZKConnect, fmt.Sprintf(util.WorldRpcService, config.Profile, config.Id), config.RpcUrl, zk.FlagEphemeral)

	log.Info("WorldManager:inited")
	return nil
}

func (this *WorldManager) Stop() {
	if this.ZKConnect != nil {
		this.ZKConnect.Close()
	}
}
