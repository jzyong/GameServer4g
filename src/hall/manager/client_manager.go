package manager

import (
	"errors"
	"fmt"
	"github.com/go-zookeeper/zk"
	"github.com/golang/protobuf/proto"
	"github.com/jzyong/go-mmo-server/src/core/log"
	network "github.com/jzyong/go-mmo-server/src/core/network/tcp"
	"github.com/jzyong/go-mmo-server/src/core/util"
	"github.com/jzyong/go-mmo-server/src/gate/handler"
	"github.com/jzyong/go-mmo-server/src/message"
	"runtime"
	"strconv"
	"sync"
)

//管理连接网关的tcp客户端
type ClientManager struct {
	util.DefaultModule
	gateClients       map[int32]*GateClient     // 网关连接
	gateClientsLock   sync.RWMutex              //网关读写锁
	MessageDistribute network.MessageDistribute //消息处理器
}

func NewClientManager() *ClientManager {
	return &ClientManager{
		gateClients:       make(map[int32]*GateClient),
		MessageDistribute: network.NewMessageDistribute(uint32(runtime.NumCPU())),
	}
}

func GetClientManager() *ClientManager {
	return Module.ClientManager
}

func (this *ClientManager) Init() error {
	log.Info("ClientManager:init")

	//开启工作线程池
	this.MessageDistribute.StartWorkerPool()

	////启动网络
	////TODO 从zookeeper中获取gate连接
	//client, err := network.NewClient("GateClient", "192.168.110.2:6061")
	//if err != nil {
	//	return err
	//}
	//this.gateClient = client
	//this.gateClient.SetChannelActive(clientChannelActive)
	//this.gateClient.SetChannelInactive(clientChannelInactive)
	//this.registerHandlers()
	//go this.gateClient.Start()
	//
	////TODO 测试发送消息 待测试接收消息处理
	//time.Sleep(time.Second * 3)
	//msg := message.UserLoginResponse{
	//	PlayerId: 1,
	//}
	//SendMsg(this.gateClient.GetChannel(), int32(message.MID_ServerListRes), 1, &msg)

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

//注册消息
func (this *ClientManager) registerHandlers() {
	this.MessageDistribute.RegisterHandler(int32(message.MID_ServerListReq), network.NewTcpHandler(handler.HandleServerList))

}

//更新网关客户端
func (this *ClientManager) UpdateGateClient(gateServerIds []string, zkConnect *zk.Conn, path string) {
	this.gateClientsLock.Lock()
	defer this.gateClientsLock.Unlock()
	//遍历添加新连接
	for _, gateIdStr := range gateServerIds {
		gateId, err := strconv.ParseInt(gateIdStr, 10, 32)
		if err != nil {
			log.Warn("网关id 异常 %v ：%v", gateIdStr, err)
		}
		if _, ok := this.gateClients[int32(gateId)]; ok {
			continue
		} else {
			//连接网关
			serverUrl := util.ZKGet(zkConnect, fmt.Sprintf("%v/%v", path, gateIdStr))
			var client = &GateClient{
				GateId:  int32(gateId),
				GateUrl: serverUrl,
			}
			this.gateClients[int32(gateId)] = client
			log.Infof("新增网关客户端：%v 地址为：%v", gateIdStr, serverUrl)
			tcpClient, err := network.NewClient(fmt.Sprintf("GateClient-%v", gateIdStr), serverUrl, this.MessageDistribute)
			if err != nil {
				log.Warn("网关id 异常 %v ：%v", gateIdStr, err)
				continue
			}
			client.Client = tcpClient
			tcpClient.SetChannelActive(clientChannelActive)
			tcpClient.SetChannelInactive(clientChannelInactive)
			go tcpClient.Start()
		}
	}
	//删除已关闭的网关
	for gateId, _ := range this.gateClients {
		gateIdStr := strconv.Itoa(int(gateId))
		if util.SliceContains(gateServerIds, gateIdStr) < 0 {
			this.gateClients[gateId].Client.Stop()
			delete(this.gateClients, gateId)
		}
	}

}

func (this *ClientManager) Stop() {
	// 关闭服务器
	for _, gateClient := range this.gateClients {
		gateClient.Client.Stop()
	}
}

//网关客户端
type GateClient struct {
	Client  network.Client //网关客户端
	GateId  int32          //网关id
	GateUrl string         //网关连接地址
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
