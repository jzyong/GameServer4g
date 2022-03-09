package main

import (
	"flag"
	"github.com/jzyong/GameServer4g/game-hall/config"
	"github.com/jzyong/GameServer4g/game-hall/handler"
	"github.com/jzyong/GameServer4g/game-hall/manager"
	"github.com/jzyong/golib/log"
	"github.com/jzyong/golib/util"
	"runtime"
)

/**
大厅入口
*/
func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())
	initConfigAndLog()

	log.Debug("hall:%d starting", config.ApplicationConfigInstance.Id)

	var err error
	err = m.Init()
	if err != nil {
		log.Error("mini game service start error: %s", err.Error())
		return
	}
	m.Run()

	//TODO 使用模块初始化？
	handler.RegisterHandlers()

	util.WaitForTerminate()
	m.Stop()
	util.WaitForTerminate()

	log.Info("hall stop")
}

//初始化项目配置和日志
func initConfigAndLog() {
	configPath := flag.String("config", "E:\\server\\GameServer4g\\game-hall\\config\\ApplicationConfig_develop.json", "配置文件加载路径")
	flag.Parse()
	config.FilePath = *configPath
	config.ApplicationConfigInstance.Reload()

	//2.关闭debug
	if "DEBUG" != config.ApplicationConfigInstance.LogLevel {
		log.CloseDebug()
	}
	log.SetLogFile("../log", "game-hall")
}

//模块管理
type ModuleManager struct {
	*util.DefaultModuleManager
	HallManager   *manager.HallManager
	ClientManager *manager.ClientManager
}

//初始化模块
func (m *ModuleManager) Init() error {
	m.HallManager = m.AppendModule(manager.GetHallManager()).(*manager.HallManager)
	m.ClientManager = m.AppendModule(manager.GetClientManager()).(*manager.ClientManager)
	return m.DefaultModuleManager.Init()
}

var m = &ModuleManager{
	DefaultModuleManager: util.NewDefaultModuleManager(),
}
