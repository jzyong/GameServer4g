package manager

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-zookeeper/zk"
	"github.com/golang/protobuf/proto"
	"github.com/jzyong/go-mmo-server/src/core/log"
	network "github.com/jzyong/go-mmo-server/src/core/network/tcp"
	"github.com/jzyong/go-mmo-server/src/core/util"
	"github.com/jzyong/go-mmo-server/src/message"
	"google.golang.org/grpc"
	"time"

	//"github.com/jzyong/go-mmo-server/src/hall/handler"
	//"github.com/jzyong/go-mmo-server/src/message"
	"runtime"
	"strconv"
	"sync"
)

////注册消息
//func (this *ClientManager) registerHandlers() {
//	//玩家
//	this.MessageDistribute.RegisterHandler(int32(message.MID_UserLoginReq), network.NewTcpHandler(handler.HandUserLogin))
//
//}

//管理连接网关的tcp客户端
type ClientManager struct {
	util.DefaultModule
	gateClients        map[int32]*GateClient     // 网关连接
	gateClientsLock    sync.RWMutex              //网关读写锁
	MessageDistribute  network.MessageDistribute //消息处理器
	WorldClientConnect *grpc.ClientConn
	PlayerWorldClient  message.PlayerWorldServiceClient
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

	//this.registerHandlers()

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
	log.Infof("创建网关连接：%v", channel.RemoteAddr())
}

//链接断开
func clientChannelInactive(channel network.Channel) {
	log.Infof("网关连接断开：%v", channel.RemoteAddr())
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

//更新world客户端
func (this *ClientManager) UpdateWorldClient(serverIds []string, zkConnect *zk.Conn, path string) {
	//遍历添加新连接
	for _, IdStr := range serverIds {
		_, err := strconv.ParseInt(IdStr, 10, 32)
		if err != nil {
			log.Warn("world id 异常 %v ：%v", IdStr, err)
		}
		//关闭之前
		if this.WorldClientConnect != nil {
			this.WorldClientConnect.Close()
		}

		//启动新连接
		serverUrl := util.ZKGet(zkConnect, fmt.Sprintf("%v/%v", path, IdStr))
		conn, err := grpc.Dial(serverUrl, grpc.WithInsecure())
		if err != nil {
			log.Warnf("%v", err)
		}
		this.WorldClientConnect = conn
		this.PlayerWorldClient = message.NewPlayerWorldServiceClient(conn)

		//TODO 测试grpc
		context, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		response, _ := this.PlayerWorldClient.Login(context, &message.UserLoginRequest{
			Account:  "player1",
			Password: "123123",
		})
		log.Infof("login return %d", response.PlayerId)

		log.Infof("connect to world %s address:%s", path, serverUrl)
	}

}

func (this *ClientManager) Stop() {
	// 关闭服务器
	for _, gateClient := range this.gateClients {
		gateClient.Client.Stop()
	}
	if this.WorldClientConnect != nil {
		this.WorldClientConnect.Close()
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

//func (this *ClientManager) GetMessageDistribute ()  network.MessageDistribute{
//	return this.MessageDistribute
//}
