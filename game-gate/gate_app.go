package main

import (
	"flag"
	"github.com/jzyong/GameServer4g/game-gate/config"
	"github.com/jzyong/GameServer4g/game-gate/handler"
	"github.com/jzyong/GameServer4g/game-gate/manager"
	"github.com/jzyong/golib/log"
	"github.com/jzyong/golib/util"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	initConfigAndLog()
	log.Debug("gate:%d starting", config.ApplicationConfigInstance.Id)

	var err error
	err = m.Init()
	if err != nil {
		log.Error("mini game service start error: %s", err.Error())
		return
	}
	m.Run()

	//TODO 使用模块初始化？
	handler.RegisterClientHandler()
	handler.RegisterGameHandler()

	util.WaitForTerminate()
	m.Stop()

	util.WaitForTerminate()

	log.Info("gate stop")
}

// 初始化项目配置和日志
func initConfigAndLog() {
	configPath := flag.String("config", "D:\\Go\\GameServer4g\\game-gate\\config\\ApplicationConfig_develop.json", "配置文件加载路径")
	flag.Parse()
	config.FilePath = *configPath
	config.ApplicationConfigInstance.Reload()

	//2.关闭debug
	if "DEBUG" != config.ApplicationConfigInstance.LogLevel {
		log.CloseDebug()
	}
	log.SetLogFile("../log", "game-gate")
}

// 模块管理
type ModuleManager struct {
	*util.DefaultModuleManager
	ClientManager *manager.ClientManager
	GameManager   *manager.GameManager
	GateManager   *manager.GateManager
	UserManager   *manager.UserManager
}

// 初始化模块
func (m *ModuleManager) Init() error {
	m.GateManager = m.AppendModule(manager.GetGateManager()).(*manager.GateManager)
	m.ClientManager = m.AppendModule(manager.GetClientManager()).(*manager.ClientManager)
	m.GameManager = m.AppendModule(manager.GetGameManager()).(*manager.GameManager)
	m.UserManager = m.AppendModule(manager.GetUserManager()).(*manager.UserManager)
	return m.DefaultModuleManager.Init()
}

var m = &ModuleManager{
	DefaultModuleManager: util.NewDefaultModuleManager(),
}
