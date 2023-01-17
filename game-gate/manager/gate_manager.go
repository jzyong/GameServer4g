package manager

import (
	"encoding/json"
	"fmt"
	"github.com/go-zookeeper/zk"
	game_common "github.com/jzyong/GameServer4g/game-common"
	"github.com/jzyong/GameServer4g/game-gate/config"
	"github.com/jzyong/golib/log"
	"github.com/jzyong/golib/util"
	"sync"
)

// 网关
type GateManager struct {
	util.DefaultModule
	ZKConnect *zk.Conn //zookeeper连接
}

var gateManager *GateManager
var gateSingletonOnce sync.Once

func GetGateManager() *GateManager {
	gateSingletonOnce.Do(func() {
		gateManager = &GateManager{}
	})
	return gateManager
}

func (m *GateManager) Init() error {
	log.Info("GateManager:init")
	//初始化id
	util.UUID = util.NewSnowflake(int16(config.ApplicationConfigInstance.Id))

	// zookeeper 初始化
	//推送配置
	configs := config.ApplicationConfigInstance
	m.ZKConnect = util.ZKCreateConnect(configs.ZookeeperUrls)
	configBytes, _ := json.Marshal(configs)
	util.ZKUpdate(m.ZKConnect, fmt.Sprintf(game_common.GateConfig, configs.Profile, configs.Id), string(configBytes))

	//注册服务
	util.ZKAdd(m.ZKConnect, fmt.Sprintf(game_common.GateGameService, configs.Profile, configs.Id), configs.GameUrl, zk.FlagEphemeral)
	util.ZKAdd(m.ZKConnect, fmt.Sprintf(game_common.GateClientService, configs.Profile, configs.Id), configs.ClientUrl, zk.FlagEphemeral)

	log.Info("GateManager:inited")
	return nil
}

func (m *GateManager) Stop() {
	if m.ZKConnect != nil {
		m.ZKConnect.Close()
	}
}
