package manager

import (
	"github.com/jzyong/GameServer4g/game-gate/config"
	"github.com/jzyong/GameServer4g/game-message/message"
	"github.com/jzyong/golib/log"
	network "github.com/jzyong/golib/network/tcp"
	"github.com/jzyong/golib/util"
	"sync"
)

//后端游戏网络管理
type GameManager struct {
	util.DefaultModule
	server        network.Server
	HallGames     map[int32]*GameServerInfo
	HallGamesLock sync.RWMutex
}

var gameManager = &GameManager{
	HallGames: make(map[int32]*GameServerInfo),
}

func GetGameManager() *GameManager {
	return gameManager
}

//
func (m *GameManager) Init() error {
	log.Info("GameManager:init")
	//启动网络
	server, err := network.NewServer("game", config.ApplicationConfigInstance.GameUrl, network.InnerServer, unregisterGameMessageDistribute)
	if err != nil {
		return err
	}
	m.server = server
	m.server.SetChannelActive(gameChannelActive)
	m.server.SetChannelInactive(gameChannelInactive)
	//m.registerHandlers()
	go m.server.Start()

	log.Info("GameManager:inited")
	return nil
}

//获取服务器
func (m *GameManager) GetServer() network.Server {
	return m.server
}

//注册消息
func (m *GameManager) RegisterHandler(mid int32, handler network.HandlerMethod) {
	m.server.RegisterHandler(mid, handler)
}

//更新服务器列表
func (m *GameManager) UpdateHallServerInfo(serverInfo *message.ServerInfo, channel network.Channel) {
	m.HallGamesLock.Lock()
	defer m.HallGamesLock.Unlock()

	hallGame, ok := m.HallGames[serverInfo.GetId()]
	if !ok {
		hallGame = &GameServerInfo{
			ServerId: serverInfo.GetId(),
		}
		m.HallGames[serverInfo.GetId()] = hallGame
		channel.SetProperty("serverId", serverInfo.GetId())
		log.Info("server %d-%d %s register to gate", serverInfo.GetType(), serverInfo.GetId(), serverInfo.GetIp())
	}
	hallGame.State = serverInfo.GetState()
	hallGame.ServerType = serverInfo.GetType()
	hallGame.Channel = channel
}

func (m *GameManager) RemoveHall(serverId int32) {
	m.HallGamesLock.Lock()
	defer m.HallGamesLock.Unlock()
	delete(m.HallGames, serverId)
}

//获取大厅后端
func (m *GameManager) GetGameServerInfo(serverId int32) *GameServerInfo {
	m.HallGamesLock.Lock()
	defer m.HallGamesLock.Unlock()
	server := m.HallGames[serverId]
	return server
}

func (m *GameManager) Stop() {
	// 关闭服务器
	if m.server != nil {
		m.server.Stop()
	}
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
		log.Info("hall server %d close", serverId)
	}
}

//转发不在本地处理的消息
func unregisterGameMessageDistribute(tcpMessage network.TcpMessage) {
	//转发给客户端
	log.Debug("转发消息：%d", tcpMessage.GetMsgId())
	user, _ := GetUserManager().GetIdUser(tcpMessage.GetObjectId())
	if user != nil {
		user.SendMessageToClient(tcpMessage)
	} else {
		log.Warn("%d send message %d fail, user not find", tcpMessage.GetObjectId(), tcpMessage.GetMsgId())
	}

}

//后端服务器
type GameServerInfo struct {
	ServerId   int32
	ServerType int32
	Channel    network.Channel
	State      int32 //服务器状态
}
