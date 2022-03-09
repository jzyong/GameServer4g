package main

import (
	"flag"
	"github.com/jzyong/go-mmo-server/src/core/log"
	"github.com/jzyong/go-mmo-server/src/core/util"
	"github.com/jzyong/go-mmo-server/src/world/config"
	"github.com/jzyong/go-mmo-server/src/world/manager"
	"runtime"
)

/**
世界服入库
*/
func main() {
	log.OpenDebug()
	log.Debug("start world")
	configPath := flag.String("config", "E:\\server\\go-mmo-server\\src\\world\\config\\WorldConfig.json", "配置文件加载路径")
	flag.Parse()
	config.FilePath = *configPath
	config.WorldConfigInstance.Reload()

	runtime.GOMAXPROCS(runtime.NumCPU())

	var err error
	err = manager.Module.Init()
	if err != nil {
		log.Errorf("world start error: %s", err.Error())
		return
	}
	manager.Module.Run()
	util.WaitForTerminate()
	manager.Module.Stop()
	log.Info("world stop")

}
