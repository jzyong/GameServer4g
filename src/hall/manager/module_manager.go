package manager

import (
	"github.com/jzyong/go-mmo-server/src/core/util"
)

type ModuleManager struct {
	*util.DefaultModuleManager
	HallManager *HallManager
}

var Module = &ModuleManager{
	DefaultModuleManager: util.NewDefaultModuleManager(),
}

func (this *ModuleManager) Init() error {
	this.HallManager = this.AppendModule(NewHallManager()).(*HallManager)
	return this.DefaultModuleManager.Init()
}
