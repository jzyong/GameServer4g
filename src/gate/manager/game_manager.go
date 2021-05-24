package manager

import (
	"github.com/jzyong/go-mmo-server/src/core/log"
	"github.com/jzyong/go-mmo-server/src/core/util"
)

//后端游戏管理
type GameManager struct {
	util.DefaultModule
}

func NewGameManager() *GameManager {
	return &GameManager{}
}

//@
func (this *GameManager) Init() error {
	log.Info("GameManager:init")
	//TODO 网络 初始化

	log.Info("GameManager:inited")
	return nil
}

func (this *GameManager) Stop() {
	//TODO 关闭服务器
	//if this.server != nil {
	//	this.server.Stop()
	//}
}
