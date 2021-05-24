package manager

import (
	"github.com/jzyong/go-mmo-server/src/core/util"
)

type ModuleManager struct {
	*util.DefaultModuleManager
	clientManager *ClientManager
	gateManager   *GateManager
	gameManager   *GameManager
	userManager   *UserManager
}

var Module = &ModuleManager{
	DefaultModuleManager: util.NewDefaultModuleManager(),
}

func (this *ModuleManager) Init() error {
	this.gateManager = this.AppendModule(NewGateManager()).(*GateManager)
	this.clientManager = this.AppendModule(NewClientManager()).(*ClientManager)
	this.gameManager = this.AppendModule(NewGameManager()).(*GameManager)
	this.userManager = this.AppendModule(NewUserManager()).(*UserManager)
	return this.DefaultModuleManager.Init()
}
