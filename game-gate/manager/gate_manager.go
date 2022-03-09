package manager

import (
	"encoding/json"
	"fmt"
	"github.com/go-zookeeper/zk"
	game_common "github.com/jzyong/GameServer4g/game-common"
	"github.com/jzyong/GameServer4g/game-gate/config"
	"github.com/jzyong/golib/log"
	"github.com/jzyong/golib/util"
)

//网关
type GateManager struct {
	util.DefaultModule
	ZKConnect *zk.Conn //zookeeper连接
}

var gateManager = &GateManager{}

func GetGateManager() *GateManager {
	return gateManager
}

//
func (this *GateManager) Init() error {
	log.Info("GateManager:init")
	//初始化id
	util.UUID = util.NewSnowflake(int16(config.ApplicationConfigInstance.Id))

	// zookeeper 初始化
	//推送配置
	config := config.ApplicationConfigInstance
	this.ZKConnect = util.ZKCreateConnect(config.ZookeeperUrls)
	configBytes, _ := json.Marshal(config)
	util.ZKUpdate(this.ZKConnect, fmt.Sprintf(game_common.GateConfig, config.Profile, config.Id), string(configBytes))

	//注册服务
	util.ZKAdd(this.ZKConnect, fmt.Sprintf(game_common.GateGameService, config.Profile, config.Id), config.GameUrl, zk.FlagEphemeral)
	util.ZKAdd(this.ZKConnect, fmt.Sprintf(game_common.GateClientService, config.Profile, config.Id), config.ClientUrl, zk.FlagEphemeral)

	log.Info("GateManager:inited")
	return nil
}

func (this *GateManager) Stop() {
	if this.ZKConnect != nil {
		this.ZKConnect.Close()
	}
}
