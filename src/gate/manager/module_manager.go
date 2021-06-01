package manager

import (
	"github.com/jzyong/go-mmo-server/src/core/util"
)

type ModuleManager struct {
	*util.DefaultModuleManager
	ClientManager *ClientManager
	GateManager   *GateManager
	GameManager   *GameManager
	UserManager   *UserManager
}

var Module = &ModuleManager{
	DefaultModuleManager: util.NewDefaultModuleManager(),
}

func (this *ModuleManager) Init() error {
	this.GateManager = this.AppendModule(NewGateManager()).(*GateManager)
	this.ClientManager = this.AppendModule(NewClientManager()).(*ClientManager)
	this.GameManager = this.AppendModule(NewGameManager()).(*GameManager)
	this.UserManager = this.AppendModule(NewUserManager()).(*UserManager)
	return this.DefaultModuleManager.Init()
}
