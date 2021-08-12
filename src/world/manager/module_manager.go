package manager

import (
	"github.com/jzyong/go-mmo-server/src/core/util"
)

type ModuleManager struct {
	*util.DefaultModuleManager
	WorldManager *WorldManager
}

var Module = &ModuleManager{
	DefaultModuleManager: util.NewDefaultModuleManager(),
}

func (this *ModuleManager) Init() error {
	this.WorldManager = this.AppendModule(NewWorldManager()).(*WorldManager)
	return this.DefaultModuleManager.Init()
}
