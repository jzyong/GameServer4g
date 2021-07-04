package manager

import (
	"encoding/json"
	"fmt"
	"github.com/go-zookeeper/zk"
	"github.com/jzyong/go-mmo-server/src/core/log"
	"github.com/jzyong/go-mmo-server/src/core/util"
	"github.com/jzyong/go-mmo-server/src/gate/config"
)

//网关
type GateManager struct {
	util.DefaultModule
	ZKConnect *zk.Conn //zookeeper连接
}

func NewGateManager() *GateManager {
	return &GateManager{}
}

//
func (this *GateManager) Init() error {
	log.Info("GateManager:init")
	//初始化id
	util.UUID = util.NewSnowflake(int16(config.GateConfigInstance.Id))

	// zookeeper 初始化
	//推送配置
	config := config.GateConfigInstance
	this.ZKConnect = util.ZKCreateConnect(config.ZookeeperUrls)
	configBytes, _ := json.Marshal(config)
	util.ZKUpdate(this.ZKConnect, fmt.Sprintf(util.GateConfig, config.Profile, config.Id), string(configBytes))

	//注册服务
	util.ZKAdd(this.ZKConnect, fmt.Sprintf(util.GateGameService, config.Profile, config.Id), config.GameUrl, zk.FlagEphemeral)
	util.ZKAdd(this.ZKConnect, fmt.Sprintf(util.GateClientService, config.Profile, config.Id), config.ClientUrl, zk.FlagEphemeral)

	log.Info("GateManager:inited")
	return nil
}

func (this *GateManager) Stop() {
	if this.ZKConnect != nil {
		this.ZKConnect.Close()
	}
}
