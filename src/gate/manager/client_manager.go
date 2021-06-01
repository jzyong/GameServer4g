package manager

import (
	"github.com/jzyong/go-mmo-server/src/core/log"
	network "github.com/jzyong/go-mmo-server/src/core/network/tcp"
	"github.com/jzyong/go-mmo-server/src/core/util"
	"github.com/jzyong/go-mmo-server/src/gate/config"
	"github.com/jzyong/go-mmo-server/src/gate/handler"
	"github.com/jzyong/go-mmo-server/src/message"
)

//客户端连接管理
type ClientManager struct {
	util.DefaultModule
	server network.Server
	//TODO 添加网络模块
}

func NewClientManager() *ClientManager {
	return &ClientManager{}
}

func GetClientManager() *ClientManager {
	return Module.ClientManager
}

func (this *ClientManager) Init() error {
	log.Info("ClientManager:init")

	//启动网络
	server, err := network.NewServer("client", config.GateConfigInstance.ClientUrl, network.ClientServer)
	if err != nil {
		return err
	}
	this.server = server
	this.server.SetChannelActive(ChannelActive)
	this.server.SetChannelInactive(ChannelInactive)
	this.registerHandlers()
	go this.server.Start()

	log.Info("ClientManager:inited")
	return nil
}

//获取服务器
func (this *ClientManager) GetServer() network.Server {
	return this.server
}

//链接激活
func ChannelActive(channel network.Channel) {
	//TODO 创建用户，加入。。。
}

//链接断开
func ChannelInactive(channel network.Channel) {
	//TODO 移除用户，。。。
}

func (this *ClientManager) registerHandlers() {
	this.server.RegisterHandler(int32(message.MID_ServerListReq), handler.HandleServerList)
}

func (this *ClientManager) Stop() {
	// 关闭服务器
	if this.server != nil {
		this.server.Stop()
	}
}
