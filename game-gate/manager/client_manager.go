package manager

import (
	"github.com/jzyong/GameServer4g/game-gate/config"
	"github.com/jzyong/golib/log"
	network "github.com/jzyong/golib/network/tcp"
	"github.com/jzyong/golib/util"
)

//客户端网络连接管理
type ClientManager struct {
	util.DefaultModule
	server network.Server
}

var clientManager = &ClientManager{}

func GetClientManager() *ClientManager {
	return clientManager
}

func (m *ClientManager) Init() error {
	log.Info("ClientManager:init")

	//启动网络
	server, err := network.NewServer("client", config.ApplicationConfigInstance.ClientUrl, network.ClientServer, unregisterMessageDistribute)
	if err != nil {
		return err
	}
	m.server = server
	m.server.SetChannelActive(clientChannelActive)
	m.server.SetChannelInactive(clientChannelInactive)
	//m.registerHandlers()
	go m.server.Start()

	log.Info("ClientManager:inited")
	return nil
}

//获取服务器
func (m *ClientManager) GetServer() network.Server {
	return m.server
}

//注册消息
func (m *ClientManager) RegisterHandler(mid int32, handler network.HandlerMethod) {
	m.server.RegisterHandler(mid, handler)
}

func (m *ClientManager) Stop() {
	// 关闭服务器
	if m.server != nil {
		m.server.Stop()
	}
}

//转发不在本地处理的消息
func unregisterMessageDistribute(tcpMessage network.TcpMessage) {
	log.Debug("转发消息：%d", tcpMessage.GetMsgId())
	u, _ := tcpMessage.GetChannel().GetProperty("user")
	user := u.(*User)
	user.SendTcpMessageToHall(tcpMessage)
}

//链接断开
func clientChannelInactive(channel network.Channel) {
	//移除用户，。。。
	u, err := channel.GetProperty("user")
	if err == nil {
		user := u.(*User)
		if user != nil {
			//log.Debug("用户退出 sessionId：", user.SessionId, " Id:", user.Id, " ip:", channel.RemoteAddr())
			GetUserManager().UserOffLine(channel, ClientClose)
		} else {
			log.Error("sessionId：%v用户不存在", channel.RemoteAddr())
		}

	} else {
		log.Warn("用户退出 ip:", channel.RemoteAddr(), " 无用户信息")
	}
}

//链接激活
func clientChannelActive(channel network.Channel) {
	// 创建用户，加入。。。
	id, _ := util.UUID.GetId()
	user := NewUser(id, channel)
	channel.SetProperty("user", user)
	GetUserManager().AddSessionUser(user)
	log.Info("用户连接创建：%v 会话：%d 总人数：%d", channel.RemoteAddr(), id, GetUserManager().GetUserCount())
}
