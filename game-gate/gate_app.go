package main

import (
	"flag"
	"github.com/jzyong/GameServer4g/game-gate/config"
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

	//TODO
	//handler.RegisterClientHandler()
	//handler.RegisterGameHandler()

	util.WaitForTerminate()
	m.Stop()

	util.WaitForTerminate()

	log.Info("gate stop")
}

//初始化项目配置和日志
func initConfigAndLog() {
	configPath := flag.String("config", "E:\\server\\go-mmo-server\\src\\gate\\config\\GateConfig.json", "配置文件加载路径")
	flag.Parse()
	config.FilePath = *configPath
	config.ApplicationConfigInstance.Reload()

	//2.关闭debug
	if "DEBUG" != config.ApplicationConfigInstance.LogLevel {
		log.CloseDebug()
	}
	log.SetLogFile("../log", "game-gate")
}

//模块管理
type ModuleManager struct {
	*util.DefaultModuleManager
	//MiniManager        *manager.MiniGameManager
	//HallClientManager  *common.HallClientManager
	//MongoClientManager *manager.DataManager
	//GrpcManager        *rpc.GRpcManager
}

//初始化模块
func (m *ModuleManager) Init() error {
	//m.MiniManager = m.AppendModule(manager.GetMiniGameManager()).(*manager.MiniGameManager)
	//m.HallClientManager = m.AppendModule(common.GetHallClientManager()).(*common.HallClientManager)
	//m.MongoClientManager = m.AppendModule(manager.GetDataManager()).(*manager.DataManager)
	//m.GrpcManager = m.AppendModule(&rpc.GRpcManager{}).(*rpc.GRpcManager)
	//TODO
	return m.DefaultModuleManager.Init()
}

var m = &ModuleManager{
	DefaultModuleManager: util.NewDefaultModuleManager(),
}
