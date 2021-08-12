package manager

import (
	"github.com/jzyong/go-mmo-server/src/core/log"
	network "github.com/jzyong/go-mmo-server/src/core/network/tcp"
	"github.com/jzyong/go-mmo-server/src/core/util"
	"github.com/jzyong/go-mmo-server/src/gate/config"
	"github.com/jzyong/go-mmo-server/src/gate/handler"
	"github.com/jzyong/go-mmo-server/src/message"
	"sync"
)

//后端游戏网络管理
type GameManager struct {
	util.DefaultModule
	server        network.Server
	HallGames     map[int32]*GameServerInfo
	HallGamesLock sync.RWMutex
}

func NewGameManager() *GameManager {
	return &GameManager{}
}

//@
func (this *GameManager) Init() error {
	log.Info("GameManager:init")
	//启动网络
	server, err := network.NewServer("game", config.GateConfigInstance.GameUrl, network.InnerServer, nil)
	if err != nil {
		return err
	}
	this.server = server
	this.server.SetChannelActive(gameChannelActive)
	this.server.SetChannelInactive(gameChannelInactive)
	this.registerHandlers()
	go this.server.Start()

	log.Info("GameManager:inited")
	return nil
}

func GetGameManager() *GameManager {
	return Module.GameManager
}

//获取服务器
func (this *GameManager) GetServer() network.Server {
	return this.server
}

//链接激活
func gameChannelActive(channel network.Channel) {
	//TODO 创建用户，加入。。。
}

//链接断开
func gameChannelInactive(channel network.Channel) {
	//TODO 移除用户，。。。
}

func (this *GameManager) registerHandlers() {
	this.server.RegisterHandler(int32(message.MID_ServerListReq), handler.HandleServerList)
}

func (this *GameManager) Stop() {
	// 关闭服务器
	if this.server != nil {
		this.server.Stop()
	}
}

//后端服务器
type GameServerInfo struct {
	ServerId   int32
	ServerType int32
	Channel    network.Channel
	State      int32 //服务器状态

}
