package manager

import (
	"encoding/json"
	"fmt"
	"github.com/go-zookeeper/zk"
	game_common "github.com/jzyong/GameServer4g/game-common"
	"github.com/jzyong/GameServer4g/game-hall/config"
	"github.com/jzyong/golib/log"
	"github.com/jzyong/golib/util"
)

//大厅
type HallManager struct {
	util.DefaultModule
	ZKConnect *zk.Conn //zookeeper连接
}

var hallManager = &HallManager{}

func GetHallManager() *HallManager {
	return hallManager
}

//
func (this *HallManager) Init() error {
	log.Info("HallManager:init")
	//初始化id
	util.UUID = util.NewSnowflake(int16(config.ApplicationConfigInstance.Id))

	// zookeeper 初始化
	//推送配置
	config := config.ApplicationConfigInstance
	this.ZKConnect = util.ZKCreateConnect(config.ZookeeperUrls)
	configBytes, _ := json.Marshal(config)
	util.ZKUpdate(this.ZKConnect, fmt.Sprintf(game_common.HallConfig, config.Profile, config.Id), string(configBytes))
	//监听网关连接
	this.watchGateService()
	this.watchWorldService()

	//注册服务
	util.ZKAdd(this.ZKConnect, fmt.Sprintf(game_common.HallRpcService, config.Profile, config.Id), config.RpcUrl, zk.FlagEphemeral)

	log.Info("HallManager:inited")
	return nil
}

//监听网关服务
func (this *HallManager) watchGateService() {
	path := fmt.Sprintf(game_common.GateGameServiceListenPath, config.ApplicationConfigInstance.Profile)
	children, errors := util.ZKWatchChildrenW(this.ZKConnect, path)
	go func() {
		for {
			select {
			case gateIds := <-children:
				log.Info("网关列表变更为：%v", gateIds)
				GetClientManager().UpdateGateClient(gateIds, this.ZKConnect, path)
			case err := <-errors:
				log.Warn("网关服务监听异常：%v", err)
			}
		}
	}()
}

//监听世界服
func (this *HallManager) watchWorldService() {
	path := fmt.Sprintf(game_common.WorldRpcServiceListenPath, config.ApplicationConfigInstance.Profile)
	children, errors := util.ZKWatchChildrenW(this.ZKConnect, path)
	go func() {
		for {
			select {
			case worldIds := <-children:
				log.Info("world变更为：%v", worldIds)
				GetClientManager().UpdateWorldClient(worldIds, this.ZKConnect, path)
			case err := <-errors:
				log.Warn("world服务监听异常：%v", err)
			}
		}
	}()
}

func (this *HallManager) Stop() {
	if this.ZKConnect != nil {
		this.ZKConnect.Close()
	}
}
