package manager

import (
	"errors"
	"fmt"
	"github.com/go-zookeeper/zk"
	"github.com/golang/protobuf/proto"
	"github.com/jzyong/GameServer4g/game-hall/config"
	"github.com/jzyong/GameServer4g/game-message/message"
	"github.com/jzyong/golib/log"
	network "github.com/jzyong/golib/network/tcp"
	"github.com/jzyong/golib/util"
	"google.golang.org/grpc"
	"runtime"
	"strconv"
	"sync"
	"time"
)

//管理连接网关的tcp客户端
type ClientManager struct {
	util.DefaultModule
	gateClients        map[int32]*GateClient     // 网关连接
	gateClientsLock    sync.RWMutex              //网关读写锁
	MessageDistribute  network.MessageDistribute //消息处理器
	WorldClientConnect *grpc.ClientConn
	PlayerWorldClient  message.PlayerWorldServiceClient
}

var clientManager = &ClientManager{
	gateClients:       make(map[int32]*GateClient),
	MessageDistribute: network.NewMessageDistribute(uint32(runtime.NumCPU()), nil),
}

func GetClientManager() *ClientManager {
	return clientManager
}

func (m *ClientManager) Init() error {
	log.Info("ClientManager:init")

	//开启工作线程池
	m.MessageDistribute.StartWorkerPool()

	//m.registerHandlers()

	////TODO 测试发送消息 待测试接收消息处理
	//time.Sleep(time.Second * 3)
	//msg := message.UserLoginResponse{
	//	PlayerId: 1,
	//}
	//SendMsg(m.gateClient.GetChannel(), int32(message.MID_ServerListRes), 1, &msg)

	//定时发送心跳
	go func() {
		for {
			for _, client := range m.gateClients {
				sendServerHeartMessage(client.Client.GetChannel())
			}
			time.Sleep(time.Second * 3)
		}
	}()

	log.Info("ClientManager:inited")
	return nil
}

//注册消息
func (m *ClientManager) RegisterHandler(mid message.MID, method network.HandlerMethod) {
	m.MessageDistribute.RegisterHandler(int32(mid), network.NewTcpHandler(method))
}

//更新网关客户端
func (m *ClientManager) UpdateGateClient(gateServerIds []string, zkConnect *zk.Conn, path string) {
	m.gateClientsLock.Lock()
	defer m.gateClientsLock.Unlock()
	//遍历添加新连接
	for _, gateIdStr := range gateServerIds {
		gateId, err := strconv.ParseInt(gateIdStr, 10, 32)
		if err != nil {
			log.Warn("网关id 异常 %v ：%v", gateIdStr, err)
		}
		if _, ok := m.gateClients[int32(gateId)]; ok {
			continue
		} else {
			//连接网关
			serverUrl := util.ZKGet(zkConnect, fmt.Sprintf("%v/%v", path, gateIdStr))
			var client = &GateClient{
				GateId:  int32(gateId),
				GateUrl: serverUrl,
			}
			m.gateClients[int32(gateId)] = client
			log.Info("新增网关客户端：%v 地址为：%v", gateIdStr, serverUrl)
			tcpClient, err := network.NewClient(fmt.Sprintf("GateClient-%v", gateIdStr), serverUrl, m.MessageDistribute)
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
	for gateId, _ := range m.gateClients {
		gateIdStr := strconv.Itoa(int(gateId))
		if util.SliceContains(gateServerIds, gateIdStr) < 0 {
			m.gateClients[gateId].Client.Stop()
			delete(m.gateClients, gateId)
		}
	}
}

//更新world客户端
func (m *ClientManager) UpdateWorldClient(serverIds []string, zkConnect *zk.Conn, path string) {
	//遍历添加新连接
	for _, IdStr := range serverIds {
		_, err := strconv.ParseInt(IdStr, 10, 32)
		if err != nil {
			log.Warn("world id 异常 %v ：%v", IdStr, err)
		}
		//关闭之前
		if m.WorldClientConnect != nil {
			m.WorldClientConnect.Close()
		}

		//启动新连接
		serverUrl := util.ZKGet(zkConnect, fmt.Sprintf("%v/%v", path, IdStr))
		conn, err := grpc.Dial(serverUrl, grpc.WithInsecure())
		if err != nil {
			log.Warn("%v", err)
		}
		m.WorldClientConnect = conn
		m.PlayerWorldClient = message.NewPlayerWorldServiceClient(conn)

		////TODO 测试grpc
		//context, cancel := context.WithTimeout(context.Background(), time.Second)
		//defer cancel()
		//response, _ := m.PlayerWorldClient.Login(context, &message.UserLoginRequest{
		//	Account:  "player1",
		//	Password: "123123",
		//})
		//log.Infof("login return %d", response.PlayerId)

		log.Info("connect to world %s address:%s", path, serverUrl)
	}

}

func (m *ClientManager) Stop() {
	// 关闭服务器
	for _, gateClient := range m.gateClients {
		gateClient.Client.Stop()
	}
	if m.WorldClientConnect != nil {
		m.WorldClientConnect.Close()
	}
}

//发送心跳消息
func sendServerHeartMessage(channel network.Channel) {
	if channel == nil {
		return
	}

	hallConfig := config.ApplicationConfigInstance
	request := &message.ServerRegisterUpdateRequest{
		ServerInfo: &message.ServerInfo{
			Id:    hallConfig.Id,
			Type:  1,
			Ip:    hallConfig.RpcUrl,
			State: 1,
		},
	}

	SendMsg(channel, int32(message.MID_ServerRegisterUpdateReq), -1, request)
}

//链接激活
func clientChannelActive(channel network.Channel) {
	// 给网关发送注册消息
	sendServerHeartMessage(channel)
	log.Info("创建网关连接：%v", channel.RemoteAddr())
	// TODO 添加属性
}

//链接断开
func clientChannelInactive(channel network.Channel) {
	log.Info("网关连接断开：%v", channel.RemoteAddr())
	//TODO 删除连接
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
