package manager

import (
	"encoding/json"
	"fmt"
	"github.com/go-zookeeper/zk"
	"github.com/jzyong/go-mmo-server/src/core/log"
	"github.com/jzyong/go-mmo-server/src/core/util"
	"github.com/jzyong/go-mmo-server/src/hall/config"
	"time"
)

//网关
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
	//监听 TODO 临时测试，监听自己 监听网关连接
	_, _, event, err := this.ZKConnect.ExistsW("/mmo/jzy/service")
	if err != nil {
		log.Errorf("zookeeper 监听失败 %v", err)
		return err
	}
	go watchZkEvent(event)
	time.Sleep(time.Second * 5)
	util.ZKUpdate(this.ZKConnect, fmt.Sprintf(util.GateConfig, config.Profile, config.Id), string(configBytes))

	//注册服务
	util.ZKAdd(this.ZKConnect, fmt.Sprintf(util.HallRpcService, config.Profile, config.Id), config.RpcUrl, zk.FlagEphemeral)

	log.Info("HallManager:inited")
	return nil
}

// zk 回调函数
func watchZkEvent(e <-chan zk.Event) {
	event := <-e
	fmt.Println("###########################")
	fmt.Println("path: ", event.Path)
	fmt.Println("type: ", event.Type.String())
	fmt.Println("state: ", event.State.String())
	fmt.Println("---------------------------")
}

func (this *HallManager) Stop() {
	if this.ZKConnect != nil {
		this.ZKConnect.Close()
	}
}
