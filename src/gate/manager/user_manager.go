package manager

import (
	"errors"
	"github.com/golang/protobuf/proto"
	"github.com/jzyong/go-mmo-server/src/core/log"
	network "github.com/jzyong/go-mmo-server/src/core/network/tcp"
	"github.com/jzyong/go-mmo-server/src/core/util"
	"github.com/jzyong/go-mmo-server/src/message"
	"sync"
)

//连接的用户管理
type UserManager struct {
	util.DefaultModule
	SessionUser     map[int64]*User //用户session ID
	SessionUserLock sync.RWMutex
	IdUser          map[int64]*User //用户角色ID
	IdUserLock      sync.RWMutex
}

func NewUserManager() *UserManager {
	return &UserManager{
		SessionUser: make(map[int64]*User),
		IdUser:      make(map[int64]*User),
	}
}

func GetUserManager() *UserManager {
	return Module.UserManager
}

//初始化
func (this *UserManager) Init() error {
	log.Info("UserManager:init")
	// 初始化

	log.Info("UserManager:inited")
	return nil
}

//向SessionMap加用户
func (manager *UserManager) AddSessionUser(user *User) {
	manager.SessionUserLock.Lock()
	defer manager.SessionUserLock.Unlock()
	//log.Debugf("%p  %p",user,*user)
	manager.SessionUser[user.SessionId] = user
}

func (manager *UserManager) AddIdUser(user *User) {
	manager.IdUserLock.Lock()
	defer manager.IdUserLock.Unlock()
	//log.Debugf("%p  %p",user,*user)
	manager.IdUser[user.Id] = user
}

func (manager *UserManager) GetSessionUser(sessionId int64) (*User, error) {
	manager.SessionUserLock.RLock()
	defer manager.SessionUserLock.RUnlock()
	if user, ok := manager.SessionUser[sessionId]; ok {
		return user, nil
	} else {
		return nil, errors.New("user not found")
	}
}

func (manager *UserManager) GetIdUser(id int64) (*User, error) {
	manager.IdUserLock.RLock()
	defer manager.IdUserLock.RUnlock()
	if user, ok := manager.IdUser[id]; ok {
		return user, nil
	} else {
		return nil, errors.New("user not found")
	}
}

func (manager *UserManager) GetUserCount() int {
	return len(manager.SessionUser)
}

func (manager *UserManager) RemoveSessionUser(sessionId int64) {
	manager.SessionUserLock.Lock()
	defer manager.SessionUserLock.Unlock()
	delete(manager.SessionUser, sessionId)
}

func (manager *UserManager) RemoveIdUser(id int64) {
	manager.IdUserLock.Lock()
	defer manager.IdUserLock.Unlock()
	delete(manager.IdUser, id)
}

//用户离线
func (manager *UserManager) UserOffLine(userChannel network.Channel, reason OffLineReason) {
	u, err := userChannel.GetProperty("user")
	if err != nil {
		log.Warnf("获取属性异常 %s", err)
		return
	}
	user := u.(*User)

	userChannel.Stop()

	user.ClientChannel = nil
	user.GameChannel = nil
	manager.RemoveIdUser(user.Id)
	manager.RemoveSessionUser(user.SessionId)
	log.Infof("%d-%v 离线因为：%d 总人数：%d", user.Id, userChannel.RemoteAddr(), reason, manager.GetUserCount())
}

func (this *UserManager) Stop() {
}

//离线原因
type OffLineReason int32

const (
	Timeout        OffLineReason = 1 // "玩家超时，服务器主动踢出"),
	DoubleLogin    OffLineReason = 2 // "玩家异地登陆或顶号，服务器断开之前的连接"),
	ClientClose    OffLineReason = 3 //, "客户端主动断开连接"),
	Exception      OffLineReason = 4 //, "服务器收到异常断开连接"),
	ServerShutdown OffLineReason = 5 //, "后端服务器关闭，断开玩家");
)

//连接的用户
type User struct {
	//角色id
	Id int64
	//回话 id
	SessionId int64
	//客户端连接
	ClientChannel network.Channel
	//游戏连接
	GameChannel network.Channel
	//登录成功
	LoginSuccess bool
	//玩家所在的游戏服
	HallId int32
	//请求消息计数
	RequestMessageCount int32
	//返回消息计数
	ResponseMessageCount int32
}

func NewUser(sessionId int64, clientChannel network.Channel) *User {
	user := &User{SessionId: sessionId, ClientChannel: clientChannel}
	return user
}

//向游戏服发消息
func (u *User) SendToHall(mid message.MID, message proto.Message) {
	if u.GameChannel == nil {
		//TODO 获取链接
		//serverInfo, _ := manager.ServerInfoManagerInstance.GetGameServerInfo(0)
		//if serverInfo == nil {
		//	log.Error("没有找到一个可用的大厅:", mid)
		//	return
		//}
		//u.GameChannel = serverInfo.Channel
	}
	network.SendProtoMsg(u.GameChannel, int32(mid), u.SessionId, u.Id, message)
}

//向游戏服发消息
func (u *User) SendTcpMessageToHall(tcpMessage network.TcpMessage) {
	if u.GameChannel == nil {
		//TODO 获取链接
		//serverInfo, _ := manager.ServerInfoManagerInstance.GetGameServerInfo(0)
		//if serverInfo == nil {
		//	log.Error("没有找到一个可用的大厅:", mid)
		//	return
		//}
		//u.GameChannel = serverInfo.Channel
	}
	network.SendMsg(u.GameChannel, tcpMessage.GetMsgId(), u.SessionId, u.Id, tcpMessage.GetData())
}

func (u *User) GetGameChanel() network.Channel {
	if u.GameChannel == nil {
		//TODO 获取链接
		//serverInfo, _ := manager.ServerInfoManagerInstance.GetGameServerInfo(0)
		//if serverInfo == nil {
		//	log.Error("没有找到一个可用的大厅:")
		//	return nil
		//}
		//u.GameChannel = serverInfo.Channel
	}
	return u.GameChannel
}

//设置服务器id
func (u *User) SetHallId(hallId int32) {
	u.HallId = hallId
}
