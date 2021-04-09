package main

import (
	"github.com/jzyong/go-mmo-server/src/core/log"
	"github.com/jzyong/go-mmo-server/src/core/util"
	"github.com/jzyong/go-mmo-server/src/gate/config"
	"github.com/jzyong/go-mmo-server/src/gate/rpc"
)

func main() {
	log.OpenDebug()
	log.Debugf("gate:%d starting", config.GateConfigInstance.Id)

	rpc.GateToClusterClient = new(rpc.GateToCluster)
	rpc.GateToClusterClient.Start(config.GateConfigInstance.ClusterRpcURL)
	rpc.RegisterToCluster()
	//select {}
	util.WaitForTerminate()
}
