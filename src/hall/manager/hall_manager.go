package manager

import (
	"encoding/json"
	"fmt"
	"github.com/go-zookeeper/zk"
	"github.com/jzyong/go-mmo-server/src/core/log"
	"github.com/jzyong/go-mmo-server/src/core/util"
	"github.com/jzyong/go-mmo-server/src/hall/config"
)

//大厅
type HallManager struct {
	util.DefaultModule
	ZKConnect *zk.Conn //zookeeper连接
}

func NewHallManager() *HallManager {
	return &HallManager{}
}

//
func (this *HallManager) Init() error {
	log.Info("HallManager:init")
	//初始化id
	util.UUID = util.NewSnowflake(int16(config.HallConfigInstance.Id))

	// zookeeper 初始化
	//推送配置
	config := config.HallConfigInstance
	this.ZKConnect = util.ZKCreateConnect(config.ZookeeperUrls)
	configBytes, _ := json.Marshal(config)
	util.ZKUpdate(this.ZKConnect, fmt.Sprintf(util.HallConfig, config.Profile, config.Id), string(configBytes))
	//监听网关连接
	this.watchGateService()

	//注册服务
	util.ZKAdd(this.ZKConnect, fmt.Sprintf(util.HallRpcService, config.Profile, config.Id), config.RpcUrl, zk.FlagEphemeral)

	log.Info("HallManager:inited")
	return nil
}

//监听网关服务
func (this *HallManager) watchGateService() {
	path := fmt.Sprintf(util.GateGameServiceListenPath, config.HallConfigInstance.Profile)
	children, errors := util.ZKWatchChildrenW(this.ZKConnect, path)
	go func() {
		for {
			select {
			case gateIds := <-children:
				log.Infof("网关列表变更为：%v", gateIds)
				GetClientManager().UpdateGateClient(gateIds, this.ZKConnect, path)
			case err := <-errors:
				log.Warnf("网关服务监听异常：%v", err)
			}
		}
	}()

}

func (this *HallManager) Stop() {
	if this.ZKConnect != nil {
		this.ZKConnect.Close()
	}
}
