package manager

import (
	"github.com/jzyong/go-mmo-server/src/core/log"
	network "github.com/jzyong/go-mmo-server/src/core/network/tcp"
	"github.com/jzyong/go-mmo-server/src/core/util"
	"github.com/jzyong/go-mmo-server/src/gate/config"
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
	return &GameManager{
		HallGames: make(map[int32]*GameServerInfo),
	}
}

//@
func (this *GameManager) Init() error {
	log.Info("GameManager:init")
	//启动网络
	server, err := network.NewServer("game", config.GateConfigInstance.GameUrl, network.InnerServer, unregisterGameMessageDistribute)
	if err != nil {
		return err
	}
	this.server = server
	this.server.SetChannelActive(gameChannelActive)
	this.server.SetChannelInactive(gameChannelInactive)
	//this.registerHandlers()
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

//注册消息
func (this GameManager) RegisterHandler(mid int32, handler network.HandlerMethod) {
	this.server.RegisterHandler(mid, handler)
}

//更新服务器列表
func (this *GameManager) UpdateHallServerInfo(serverInfo *message.ServerInfo, channel network.Channel) {
	this.HallGamesLock.Lock()
	defer this.HallGamesLock.Unlock()

	hallGame, ok := this.HallGames[serverInfo.GetId()]
	if !ok {
		hallGame = &GameServerInfo{
			ServerId: serverInfo.GetId(),
		}
		this.HallGames[serverInfo.GetId()] = hallGame
		channel.SetProperty("serverId", serverInfo.GetId())
		log.Infof("server %d-%d %s register to gate", serverInfo.GetType(), serverInfo.GetId(), serverInfo.GetIp())
	}
	hallGame.State = serverInfo.GetState()
	hallGame.ServerType = serverInfo.GetType()
	hallGame.Channel = channel
}

func (this *GameManager) RemoveHall(serverId int32) {
	this.HallGamesLock.Lock()
	defer this.HallGamesLock.Unlock()
	delete(this.HallGames, serverId)
}

//获取大厅后端
func (this *GameManager) GetGameServerInfo(serverId int32) *GameServerInfo {
	this.HallGamesLock.Lock()
	defer this.HallGamesLock.Unlock()
	server := this.HallGames[serverId]
	return server
}

//链接激活
func gameChannelActive(channel network.Channel) {
	//TODO 创建用户，加入。。。
}

//链接断开
func gameChannelInactive(channel network.Channel) {
	// 移除服务器
	id, err := channel.GetProperty("serverId")
	if err == nil {
		serverId := id.(int32)
		GetGameManager().RemoveHall(serverId)
		log.Infof("hall server %d close", serverId)
	}
}

//转发不在本地处理的消息
func unregisterGameMessageDistribute(tcpMessage network.TcpMessage) {
	//转发给客户端
	log.Debugf("转发消息：%d", tcpMessage.GetMsgId())
	user, _ := GetUserManager().GetIdUser(tcpMessage.GetObjectId())
	if user != nil {
		user.SendMessageToClient(tcpMessage)
	} else {
		log.Warnf("%d send message %d fail, user not find", tcpMessage.GetObjectId(), tcpMessage.GetMsgId())
	}

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
