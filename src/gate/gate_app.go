package main

import (
	"flag"
	"github.com/jzyong/go-mmo-server/src/core/log"
	"github.com/jzyong/go-mmo-server/src/core/util"
	"github.com/jzyong/go-mmo-server/src/gate/config"
	"github.com/jzyong/go-mmo-server/src/gate/handler"
	"github.com/jzyong/go-mmo-server/src/gate/manager"
	"runtime"
)

func main() {
	log.OpenDebug()
	//log.SetLogFile("../../log","gate.log") //正式服需要输出到文件
	log.Debugf("gate:%d starting", config.GateConfigInstance.Id)

	configPath := flag.String("config", "E:\\server\\go-mmo-server\\src\\gate\\config\\GateConfig.json", "配置文件加载路径")
	flag.Parse()
	config.FilePath = *configPath
	config.GateConfigInstance.Reload()

	runtime.GOMAXPROCS(runtime.NumCPU())

	var err error
	err = manager.Module.Init()
	if err != nil {
		log.Errorf("gate start error: %s", err.Error())
		return
	}
	manager.Module.Run()

	handler.RegisterClientHandler()
	handler.RegisterGameHandler()
	util.WaitForTerminate()
	manager.Module.Stop()
	log.Info("gate stop")
}
