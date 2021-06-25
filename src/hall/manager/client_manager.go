package manager

import (
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/jzyong/go-mmo-server/src/core/log"
	network "github.com/jzyong/go-mmo-server/src/core/network/tcp"
	"github.com/jzyong/go-mmo-server/src/core/util"
	"github.com/jzyong/go-mmo-server/src/gate/handler"
	"github.com/jzyong/go-mmo-server/src/message"
	"time"
)

//管理连接网关的tcp客户端
type ClientManager struct {
	util.DefaultModule
	gateClient network.Client //网关客户端 TODO 修改为列表
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
	//TODO 从zookeeper中获取gate连接
	client, err := network.NewClient("GateClient", "192.168.110.2:6061")
	if err != nil {
		return err
	}
	this.gateClient = client
	this.gateClient.SetChannelActive(clientChannelActive)
	this.gateClient.SetChannelInactive(clientChannelInactive)
	this.registerHandlers()
	go this.gateClient.Start()

	//TODO 测试发送消息 待测试接收消息处理
	time.Sleep(time.Second * 3)
	msg := message.UserLoginResponse{
		PlayerId: 1,
	}
	SendMsg(this.gateClient.GetChannel(), int32(message.MID_ServerListRes), 1, &msg)

	log.Info("ClientManager:inited")
	return nil
}

//链接激活
func clientChannelActive(channel network.Channel) {
	//TODO
	//// 创建用户，加入。。。
	//id, _ := util.UUID.GetId()
	//user := NewUser(id, channel)
	//channel.SetProperty("user", user)
	//GetUserManager().AddSessionUser(user)
	//log.Infof("用户连接创建：%v 会话：%d 总人数：%d", channel.RemoteAddr(), id, GetUserManager().GetUserCount())
}

//链接断开
func clientChannelInactive(channel network.Channel) {
	//TODO
	////移除用户，。。。
	//u, err := channel.GetProperty("user")
	//if err == nil {
	//	user := u.(*User)
	//	if user != nil {
	//		//log.Debug("用户退出 sessionId：", user.SessionId, " Id:", user.Id, " ip:", channel.RemoteAddr())
	//		GetUserManager().UserOffLine(channel, ClientClose)
	//	} else {
	//		log.Errorf("sessionId：%v用户不存在", channel.RemoteAddr())
	//	}
	//
	//} else {
	//	log.Warn("用户退出 ip:", channel.RemoteAddr(), " 无用户信息")
	//}
}

func (this *ClientManager) registerHandlers() {
	this.gateClient.RegisterHandler(int32(message.MID_ServerListReq), handler.HandleServerList)
}

//发送消息
func SendMsg(channel network.Channel, msgId int32, senderId int64, message proto.Message) error {
	data, err := proto.Marshal(message)
	if err != nil {
		fmt.Println("Pack error msg id = ", msgId)
		return errors.New("Pack error msg ")
	}
	if channel.IsClose() == true {
		return errors.New("connection closed when send msg")
	}
	//将data封包，并且发送
	var decoder = network.NewInnerDataPack()
	msg, err := decoder.Pack(network.NewInnerMessage(msgId, data, senderId, 0))
	if err != nil {
		fmt.Println("Pack error msg id = ", msgId)
		return errors.New("Pack error msg ")
	}
	//写回客户端
	channel.GetMsgChan() <- msg
	return nil
}

func (this *ClientManager) Stop() {
	// 关闭服务器
	if this.gateClient != nil {
		this.gateClient.Stop()
	}
}
