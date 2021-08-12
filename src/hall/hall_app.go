package main

import (
	"flag"
	"github.com/jzyong/go-mmo-server/src/core/log"
	"github.com/jzyong/go-mmo-server/src/core/util"
	"github.com/jzyong/go-mmo-server/src/hall/config"
	"github.com/jzyong/go-mmo-server/src/hall/manager"
	"runtime"
)

/**
大厅入口
*/
func main() {
	log.OpenDebug()
	//log.SetLogFile("../../log","gate.log") //正式服需要输出到文件
	log.Debugf("hall:%d starting", config.HallConfigInstance.Id)

	configPath := flag.String("config", "E:\\server\\go-mmo-server\\src\\hall\\config\\HallConfig.json", "配置文件加载路径")
	flag.Parse()
	config.FilePath = *configPath
	config.HallConfigInstance.Reload()

	runtime.GOMAXPROCS(runtime.NumCPU())

	var err error
	err = manager.Module.Init()
	if err != nil {
		log.Errorf("hall start error: %s", err.Error())
		return
	}
	manager.Module.Run()
	util.WaitForTerminate()
	manager.Module.Stop()
	log.Info("hall stop")
}
