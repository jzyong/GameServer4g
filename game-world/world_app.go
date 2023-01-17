package main

import (
	"flag"
	"github.com/jzyong/GameServer4g/game-world/config"
	"github.com/jzyong/GameServer4g/game-world/manager"
	"github.com/jzyong/golib/log"
	"github.com/jzyong/golib/util"
	"runtime"
)

/*
*
世界服入库
*/
func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	initConfigAndLog()
	log.Debug("start world")

	var err error
	err = m.Init()
	if err != nil {
		log.Error("mini game service start error: %s", err.Error())
		return
	}
	m.Run()

	util.WaitForTerminate()
	m.Stop()
	util.WaitForTerminate()
	log.Info("world stop")

}

// 初始化项目配置和日志
func initConfigAndLog() {
	configPath := flag.String("config", "D:\\Go\\GameServer4g\\game-world\\config\\ApplicationConfig_develop.json", "配置文件加载路径")
	flag.Parse()
	config.FilePath = *configPath
	config.ApplicationConfigInstance.Reload()

	//2.关闭debug
	if "DEBUG" != config.ApplicationConfigInstance.LogLevel {
		log.CloseDebug()
	}
	log.SetLogFile("../log", "game-hall")
}

// 模块管理
type ModuleManager struct {
	*util.DefaultModuleManager
	WorldManager *manager.WorldManager
}

// 初始化模块
func (m *ModuleManager) Init() error {
	m.WorldManager = m.AppendModule(manager.GetWorldManager()).(*manager.WorldManager)
	return m.DefaultModuleManager.Init()
}

var m = &ModuleManager{
	DefaultModuleManager: util.NewDefaultModuleManager(),
}
