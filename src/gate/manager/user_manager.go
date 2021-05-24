package manager

import (
	"github.com/jzyong/go-mmo-server/src/core/log"
	"github.com/jzyong/go-mmo-server/src/core/util"
)

//连接的用户管理
type UserManager struct {
	util.DefaultModule
}

func NewUserManager() *UserManager {
	return &UserManager{}
}

//
func (this *UserManager) Init() error {
	log.Info("UserManager:init")
	//TODO 初始化

	log.Info("UserManager:inited")
	return nil
}

func (this *UserManager) Stop() {
	//TODO 关闭服务器
	//if this.server != nil {
	//	this.server.Stop()
	//}
}
