package network

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/jzyong/go-mmo-server/src/core/log"
	"io"
	"net"
	"runtime"
	"sync"
)

//定义服务器接口
type Client interface {
	//启动客户端方法
	Start()
	//停止客户端方法
	Stop()
	//开启业务服务方法
	Run()
	//路由功能：给当前服务注册一个路由业务方法，供客户端链接处理使用
	RegisterHandler(msgId int32, handler HandlerMethod)
	//得到链接
	GetChannel() Channel
	//设置该Client的连接创建时Hook函数
	SetChannelActive(func(Channel))
	//设置该Client的连接断开时的Hook函数
	SetChannelInactive(func(Channel))
	//调用连接OnConnStart Hook函数
	ChannelActive(conn Channel)
	//调用连接OnConnStop Hook函数
	ChannelInactive(conn Channel)
}

//Client 接口实现，定义一个Server服务类
type clientImpl struct {
	//服务器的名称
	Name string
	//服务绑定的地址
	ServerUrl string
	//当前Server的消息管理模块，用来绑定MsgId和对应的处理方法
	MessageDistribute MessageDistribute
	//当前Server的链接管理器
	Channel Channel
	//该Client的连接创建时Hook函数
	channelActive func(conn Channel)
	//该Client的连接断开时的Hook函数
	channelInactive func(conn Channel)
	//客户端连接
	Conn net.Conn
}

//创建网络服务
func NewClient(name, url string) (Client, error) {
	return &clientImpl{
		Name:              name,
		ServerUrl:         url,
		MessageDistribute: NewMessageDistribute(uint32(runtime.NumCPU())),
	}, nil
}

//============== 实现 Client 里的全部接口方法 ========

//开启网络服务 用go启动
func (s *clientImpl) Start() {
	log.Infof("[START] Client name: %s,connect to %s is starting\n", s.Name, s.ServerUrl)

	//开启一个go去做服务端Linster业务
	go func() {
		//0 启动worker工作池机制
		s.MessageDistribute.StartWorkerPool()

		//1 连接服务器地址
		conn, err := net.Dial("tcp", s.ServerUrl)
		if err != nil {
			log.Errorf("Game start err, exit! %v", err)
			return
		}
		//2 已经监听成功
		log.Info("client ", s.Name, " success, now connecting...")
		s.Conn = conn
		channel := NewClientChannel(conn, s.MessageDistribute, s)

		//3 启动当前链接的处理业务
		go channel.Start()
	}()
	//阻塞,否则主Go退出， listenner的go将会退出
	select {}
}

//停止服务
func (s *clientImpl) Stop() {
	log.Infof("客户端%s连接关闭", s.ServerUrl)
	s.Conn.Close()
}

//运行服务
func (s *clientImpl) Run() {
	s.Start()
	//阻塞,否则主Go退出， listener的go将会退出
	select {}
}

//路由功能：给当前服务注册一个路由业务方法，供客户端链接处理使用
func (s *clientImpl) RegisterHandler(msgId int32, handler HandlerMethod) {
	s.MessageDistribute.RegisterHandler(msgId, NewTcpHandler(handler))
}

//得到链接
func (s *clientImpl) GetChannel() Channel {
	return s.Channel
}

//设置该Server的连接创建时Hook函数
func (s *clientImpl) SetChannelActive(hookFunc func(Channel)) {
	s.channelActive = hookFunc
}

//设置该Server的连接断开时的Hook函数
func (s *clientImpl) SetChannelInactive(hookFunc func(Channel)) {
	s.channelInactive = hookFunc
}

//调用连接OnConnStart Hook函数
func (s *clientImpl) ChannelActive(conn Channel) {
	if s.channelActive != nil {
		//fmt.Println("---> CallOnConnStart....")
		s.channelActive(conn)
	}
}

//调用连接OnConnStop Hook函数
func (s *clientImpl) ChannelInactive(conn Channel) {
	if s.channelInactive != nil {
		//fmt.Println("---> CallOnConnStop....")
		s.channelInactive(conn)
	}
}

//连接会话 实现Channel
type clientChannelImpl struct {
	Channel
	//当前连接的连接
	Conn net.Conn
	//当前连接的关闭状态
	IsClosed bool
	//消息管理MsgId和对应处理方法的消息管理模块
	MessageDistribute MessageDistribute
	//告知该链接已经退出/停止的channel
	ExitBuffChan chan bool
	//无缓冲管道，用于读、写两个goroutine之间的消息通信
	msgChan chan []byte
	//有缓冲管道，用于读、写两个goroutine之间的消息通信
	msgBuffChan chan []byte
	//链接属性
	property map[string]interface{}
	//保护链接属性修改的锁
	propertyLock sync.RWMutex
	//连接的客户端
	Client Client
}

//创建连接的方法
func NewClientChannel(conn net.Conn, messageDistribute MessageDistribute, client Client) Channel {
	//初始化Conn属性
	c := &clientChannelImpl{
		Conn:              conn,
		IsClosed:          false,
		MessageDistribute: messageDistribute,
		ExitBuffChan:      make(chan bool, 1),
		msgChan:           make(chan []byte),
		msgBuffChan:       make(chan []byte, 1024),
		property:          make(map[string]interface{}),
		Client:            client,
	}
	return c
}

//	写消息Goroutine， 用户将数据发送给客户端
func (c *clientChannelImpl) StartWriter() {
	defer log.Infof("%s conn Writer exit", c.RemoteAddr().String())

	for {
		select {
		case data := <-c.msgChan:
			//有数据要写给客户端
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send Data error:, ", err, " Conn Writer exit")
				return
			}
			//fmt.Printf("Send data succ! data = %+v\n", data)
		case data, ok := <-c.msgBuffChan:
			if ok {
				//有数据要写给客户端
				if _, err := c.Conn.Write(data); err != nil {
					fmt.Println("Send Buff Data error:, ", err, " Conn Writer exit")
					return
				}
			} else {
				fmt.Println("msgBuffChan is Closed")
				break
			}
		case <-c.ExitBuffChan:
			return
		}
	}
}

//	读消息Goroutine，用于从客户端中读取数据
func (c *clientChannelImpl) StartReader() {
	defer log.Infof("%s conn Reader exit!", c.RemoteAddr().String())
	defer c.Stop()
	for {
		// 创建拆包解包的对象
		buffMsgLength := make([]byte, 4)

		// read len
		var decoder = NewInnerDataPack()

		if _, err := io.ReadFull(c.Conn, buffMsgLength); err != nil {
			fmt.Println("read msg length error", err)
		}
		var msgLength = uint32(binary.LittleEndian.Uint32(buffMsgLength))
		// 最大长度验证
		if msgLength > 10000 {
			log.Warnf("消息太长：%d\n", msgLength)
		}

		msgData := make([]byte, msgLength)

		if _, err := io.ReadFull(c.Conn, msgData); err != nil {
			fmt.Println("read msg data error ", err)
			break
		}
		//拆包，得到msgid 和 数据 放在msg中
		msg, err := decoder.Unpack(msgData, msgLength)
		if err != nil {
			log.Error("unpack error ", err)
			break
		}
		msg.SetChannel(c)
		c.InboundHandler(msg)
	}
}

//子类实现
func (c *clientChannelImpl) InboundHandler(msg TcpMessage) {
	c.MessageDistribute.RunHandler(msg)
}

//启动连接，让当前连接开始工作
func (c *clientChannelImpl) Start() {
	//1 开启用户从客户端读取数据流程的Goroutine
	go c.StartReader()
	//2 开启用于写回客户端数据流程的Goroutine
	go c.StartWriter()
	//按照用户传递进来的创建连接时需要处理的业务，执行钩子方法
	c.Client.ChannelActive(c)
}

//停止连接，结束当前连接状态M
func (c *clientChannelImpl) Stop() {
	log.Infof("Conn Stop()...")
	//如果当前链接已经关闭
	if c.IsClosed == true {
		return
	}
	c.IsClosed = true

	//如果用户注册了该链接的关闭回调业务，那么在此刻应该显示调用
	c.Client.ChannelInactive(c)

	// 关闭socket链接
	c.Conn.Close()
	//关闭Writer
	c.ExitBuffChan <- true

	//关闭该链接全部管道
	close(c.ExitBuffChan)
	close(c.msgBuffChan)
}

////从当前连接获取原始的socket TCPConn
//func (c *clientChannelImpl) GetTCPConnection() *net.TCPConn {
//	return nil
//}
//
////获取当前连接ID
//func (c *clientChannelImpl) GetConnID() uint32 {
//	return 1
//}

//获取远程客户端地址信息
func (c *clientChannelImpl) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

//直接将Message数据发送数据给远程的TCP客户端
func (c *clientChannelImpl) SendMsg(message TcpMessage) error {
	if c.IsClosed == true {
		return errors.New("connection closed when send msg")
	}
	//将data封包，并且发送
	var decoder = NewInnerDataPack()
	msg, err := decoder.Pack(message)
	if err != nil {
		fmt.Println("Pack error msg id = ", message.GetMsgId())
		return errors.New("Pack error msg ")
	}
	//写回客户端
	c.msgChan <- msg

	return nil
}

//是否关闭
func (c *clientChannelImpl) IsClose() bool {
	return c.IsClosed
}

func (c *clientChannelImpl) GetServer() Server {
	return nil
}

func (c *clientChannelImpl) GetMsgChan() chan []byte {
	return c.msgChan
}

func (c *clientChannelImpl) GetMsgBuffChan() chan []byte {
	return c.msgBuffChan
}

//设置链接属性
func (c *clientChannelImpl) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	c.property[key] = value
}

//获取链接属性
func (c *clientChannelImpl) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()
	if value, ok := c.property[key]; ok {
		return value, nil
	} else {
		return nil, errors.New("no property found")
	}
}

//移除链接属性
func (c *clientChannelImpl) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	delete(c.property, key)
}
