package manager

import (
	"github.com/jzyong/go-mmo-server/src/core/util"
)

type ModuleManager struct {
	*util.DefaultModuleManager
	serverSeq int

	//TODO 添加其他Manager引用
	//gateManager  *GateManager
	//gsManager    *GSManager
	//fightManager *fight.FightManager
}

var Module = &ModuleManager{
	DefaultModuleManager: util.NewDefaultModuleManager(),
}

func (this *ModuleManager) Init() error {
	//this.gateManager = this.AppendModule(NewGateManager()).(*GateManager)
	//this.gsManager = this.AppendModule(NewGSManager()).(*GSManager)
	//this.fightManager = this.AppendModule(fight.NewFightManager()).(*fight.FightManager)
	//TODO
	return this.DefaultModuleManager.Init()
}
