package manager

import (
	"github.com/jzyong/go-mmo-server/src/core/log"
	"github.com/jzyong/go-mmo-server/src/core/util"
)

//客户端连接管理
type ClientManager struct {
	util.DefaultModule
	//TODO 添加网络模块
}

func NewClientManager() *ClientManager {
	return &ClientManager{}
}

//@
func (this *ClientManager) Init() error {
	log.Info("ClientManager:init")
	//TODO 添加网络模块
	//context := &nw.Context{
	//	SessionCreator: func(conn nw.Conn) nw.Session { return NewClientSession(conn) },
	//	Splitter:       pb.Split,
	//	ChanSize:       200,
	//}
	//server := wsserver.NewServer(context)
	//err := server.Start(conf.GetPort())
	//if err != nil {
	//	return err
	//}
	//this.server = server

	log.Info("ClientManager:inited")
	return nil
}

func (this *ClientManager) Stop() {
	//TODO 关闭服务器
	//if this.server != nil {
	//	this.server.Stop()
	//}
}
