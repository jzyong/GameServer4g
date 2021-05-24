package manager

import (
	"github.com/jzyong/go-mmo-server/src/core/log"
	"github.com/jzyong/go-mmo-server/src/core/util"
)

//网关
type GateManager struct {
	util.DefaultModule
}

func NewGateManager() *GateManager {
	return &GateManager{}
}

//@
func (this *GateManager) Init() error {
	log.Info("GateManager:init")
	//TODO zookeeper 初始化

	log.Info("GateManager:inited")
	return nil
}

func (this *GateManager) Stop() {
	//TODO 关闭服务器
	//if this.server != nil {
	//	this.server.Stop()
	//}
}
