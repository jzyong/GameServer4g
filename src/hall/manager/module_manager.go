package manager

import (
	"github.com/jzyong/go-mmo-server/src/core/util"
)

type ModuleManager struct {
	*util.DefaultModuleManager
	HallManager   *HallManager
	ClientManager *ClientManager
}

var Module = &ModuleManager{
	DefaultModuleManager: util.NewDefaultModuleManager(),
}

func (this *ModuleManager) Init() error {
	this.HallManager = this.AppendModule(NewHallManager()).(*HallManager)
	this.ClientManager = this.AppendModule(NewClientManager()).(*ClientManager)
	return this.DefaultModuleManager.Init()
}
