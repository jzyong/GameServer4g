package network

import (
	"fmt"
	"github.com/jzyong/go-mmo-server/src/core/log"
	"net"
	"reflect"
	"runtime"
	"strconv"
	"strings"
)

const (
	ClientServer int32 = 0
	InnerServer  int32 = 1
)

//定义服务器接口
type Server interface {
	//启动服务器方法
	Start()
	//停止服务器方法
	Stop()
	//开启业务服务方法
	Run()
	//路由功能：给当前服务注册一个路由业务方法，供客户端链接处理使用
	RegisterHandler(msgId int32, handler TcpHandler)
	//得到链接管理
	GetChannelManager() ChannelManager
	//设置该Server的连接创建时Hook函数
	SetChannelActive(func(Channel))
	//设置该Server的连接断开时的Hook函数
	SetChannelInactive(func(Channel))
	//调用连接OnConnStart Hook函数
	ChannelActive(conn Channel)
	//调用连接OnConnStop Hook函数
	ChannelInactive(conn Channel)
	//服务器类型 0客户端，1内部服务器
	GetServerType() int32
}

//iServer 接口实现，定义一个Server服务类
type serverImpl struct {
	//服务器的名称
	Name string
	//tcp4 or other
	IPVersion string
	//服务绑定的IP地址
	IP string
	//服务绑定的端口
	Port int32
	//当前Server的消息管理模块，用来绑定MsgId和对应的处理方法
	MessageDistribute MessageDistribute
	//当前Server的链接管理器
	ChannelManager ChannelManager
	//该Server的连接创建时Hook函数
	channelActive func(conn Channel)
	//该Server的连接断开时的Hook函数
	channelInactive func(conn Channel)
	///服务器类型
	ServerType int32
}

//创建网络服务
func NewServer(name, url string, serverType int32) (Server, error) {
	ipPorts := strings.Split(url, ":")
	if len(ipPorts) < 2 {
		return nil, fmt.Errorf("url format error")
	}
	port, err := strconv.ParseInt(ipPorts[1], 10, 32)
	if err != nil {
		return nil, err
	}

	return &serverImpl{
		Name:              name,
		IPVersion:         "tcp4", //暂时都是tcp4
		IP:                ipPorts[0],
		Port:              int32(port),
		MessageDistribute: NewMessageDistribute(uint32(runtime.NumCPU())),
		ChannelManager:    NewChannelManager(),
		ServerType:        serverType,
	}, nil

}

//============== 实现 Server 里的全部接口方法 ========

//开启网络服务 用go启动
func (s *serverImpl) Start() {
	log.Infof("[START] Server name: %s,listener at IP: %s, Port %d is starting\n", s.Name, s.IP, s.Port)

	//开启一个go去做服务端Linster业务
	go func() {
		//0 启动worker工作池机制
		s.MessageDistribute.StartWorkerPool()

		//1 获取一个TCP的Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr err: ", err)
			return
		}

		//2 监听服务器地址
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen", s.IPVersion, "err", err)
			return
		}

		//已经监听成功
		log.Info("start server  ", s.Name, " success, now listening...")

		//TODO server.go 应该有一个自动生成ID的方法
		var cid uint32
		cid = 0

		//3 启动server网络连接业务
		for {
			//3.1 阻塞等待客户端建立连接请求
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err ", err)
				continue
			}
			log.Info("Get conn remote addr = ", conn.RemoteAddr().String())

			//3.2 设置服务器最大连接控制,如果超过最大连接，那么则关闭此新的连接
			if s.ChannelManager.Len() >= 10000 {
				conn.Close()
				continue
			}

			//3.3 处理该新连接请求的 业务 方法， 此时应该有 handler 和 conn是绑定的
			//dealConn := NewChannel(s, conn, cid, s.MessageDistribute)
			channel := s.NewChannel(conn, cid)
			cid++

			//3.4 启动当前链接的处理业务
			go channel.Start()
		}
	}()
	//阻塞,否则主Go退出， listenner的go将会退出
	select {}
}

func (s *serverImpl) NewChannel(conn *net.TCPConn, cid uint32) Channel {
	c := NewChannel(s, conn, cid, s.MessageDistribute)
	return c
}

//停止服务
func (s *serverImpl) Stop() {
	fmt.Println("[STOP] Zinx server , name ", s.Name)
	//将其他需要清理的连接信息或者其他信息 也要一并停止或者清理
	s.ChannelManager.ClearConn()
}

//运行服务
func (s *serverImpl) Run() {
	s.Start()
	//阻塞,否则主Go退出， listener的go将会退出
	select {}
}

//路由功能：给当前服务注册一个路由业务方法，供客户端链接处理使用
func (s *serverImpl) RegisterHandler(msgId int32, handler TcpHandler) {
	log.Infof("%d 注册handler %s", handler, reflect.TypeOf(handler).Elem().Name())
	s.MessageDistribute.RegisterHandler(msgId, handler)
}

//得到链接管理
func (s *serverImpl) GetChannelManager() ChannelManager {
	return s.ChannelManager
}

//设置该Server的连接创建时Hook函数
func (s *serverImpl) SetChannelActive(hookFunc func(Channel)) {
	s.channelActive = hookFunc
}

//设置该Server的连接断开时的Hook函数
func (s *serverImpl) SetChannelInactive(hookFunc func(Channel)) {
	s.channelInactive = hookFunc
}

//调用连接OnConnStart Hook函数
func (s *serverImpl) ChannelActive(conn Channel) {
	if s.ChannelActive != nil {
		//fmt.Println("---> CallOnConnStart....")
		s.ChannelActive(conn)
	}
}

//调用连接OnConnStop Hook函数
func (s *serverImpl) ChannelInactive(conn Channel) {
	if s.ChannelInactive != nil {
		//fmt.Println("---> CallOnConnStop....")
		s.ChannelInactive(conn)
	}
}

//服务器类型
func (s serverImpl) GetServerType() int32 {
	return s.ServerType
}
