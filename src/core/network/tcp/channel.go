package network

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/jzyong/go-mmo-server/src/core/log"
	"io"
	"net"
	"sync"
)

//定义连接接口
type Channel interface {
	//启动连接，让当前连接开始工作
	Start()
	//停止连接，结束当前连接状态M
	Stop()
	//从当前连接获取原始的socket TCPConn
	GetTCPConnection() *net.TCPConn
	//获取当前连接ID
	GetConnID() uint32
	//获取远程客户端地址信息
	RemoteAddr() net.Addr
	//发送消息
	SendMsg(message TcpMessage) error
	//连接是否关闭
	IsClose() bool
	//关联的Server
	GetServer() Server
	//设置链接属性
	SetProperty(key string, value interface{})
	//获取链接属性
	GetProperty(key string) (interface{}, error)
	//移除链接属性
	RemoveProperty(key string)
	//无缓冲管道，用于读、写两个goroutine之间的消息通信
	GetMsgChan() chan []byte
	//有缓冲管道，用于读、写两个goroutine之间的消息通信
	GetMsgBuffChan() chan []byte
}

//连接会话 实现Channel
type channelImpl struct {
	//当前Conn属于哪个Server
	TcpServer Server
	//当前连接的socket TCP套接字
	Conn *net.TCPConn
	//当前连接的ID 也可以称作为SessionID，ID全局唯一
	ConnID uint32
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
}

//创建连接的方法
func NewChannel(server Server, conn *net.TCPConn, connID uint32, messageDistribute MessageDistribute) Channel {
	//初始化Conn属性
	c := &channelImpl{
		TcpServer:         server,
		Conn:              conn,
		ConnID:            connID,
		IsClosed:          false,
		MessageDistribute: messageDistribute,
		ExitBuffChan:      make(chan bool, 1),
		msgChan:           make(chan []byte),
		msgBuffChan:       make(chan []byte, 1024),
		property:          make(map[string]interface{}),
	}
	//将新创建的Conn添加到链接管理中
	c.TcpServer.GetChannelManager().Add(c)
	return c
}

//	写消息Goroutine， 用户将数据发送给客户端
func (c *channelImpl) StartWriter() {
	//fmt.Println("[Writer Goroutine is running]")
	defer fmt.Println(c.RemoteAddr().String(), "[conn Writer exit!]")

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
func (c *channelImpl) StartReader() {
	//fmt.Println("[Reader Goroutine is running]")
	defer fmt.Println(c.RemoteAddr().String(), "[conn Reader exit!]")
	defer c.Stop()

	for {
		// 创建拆包解包的对象
		buffMsgLength := make([]byte, 4)

		// read len
		var decoder DataPack
		if c.TcpServer.GetServerType() == ClientServer {
			decoder = NewClientDataPack()
		} else {
			decoder = NewInnerDataPack()
		}

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
func (c *channelImpl) InboundHandler(msg TcpMessage) {
	c.MessageDistribute.RunHandler(msg)
}

//启动连接，让当前连接开始工作
func (c *channelImpl) Start() {
	//1 开启用户从客户端读取数据流程的Goroutine
	go c.StartReader()
	//2 开启用于写回客户端数据流程的Goroutine
	go c.StartWriter()
	//按照用户传递进来的创建连接时需要处理的业务，执行钩子方法
	c.TcpServer.OnConnStart(c)
}

//停止连接，结束当前连接状态M
func (c *channelImpl) Stop() {
	fmt.Println("Conn Stop()...ConnID = ", c.ConnID)
	//如果当前链接已经关闭
	if c.IsClosed == true {
		return
	}
	c.IsClosed = true

	//如果用户注册了该链接的关闭回调业务，那么在此刻应该显示调用
	c.TcpServer.OnConnStop(c)

	// 关闭socket链接
	c.Conn.Close()
	//关闭Writer
	c.ExitBuffChan <- true

	//将链接从连接管理器中删除
	c.TcpServer.GetChannelManager().Remove(c)

	//关闭该链接全部管道
	close(c.ExitBuffChan)
	close(c.msgBuffChan)
}

//从当前连接获取原始的socket TCPConn
func (c *channelImpl) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

//获取当前连接ID
func (c *channelImpl) GetConnID() uint32 {
	return c.ConnID
}

//获取远程客户端地址信息
func (c *channelImpl) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

//直接将Message数据发送数据给远程的TCP客户端
func (c *channelImpl) SendMsg(message TcpMessage) error {
	if c.IsClosed == true {
		return errors.New("connection closed when send msg")
	}
	//将data封包，并且发送
	var decoder DataPack
	if c.TcpServer.GetServerType() == ClientServer {
		decoder = NewClientDataPack()
	} else {
		decoder = NewInnerDataPack()
	}
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
func (c *channelImpl) IsClose() bool {
	return c.IsClosed
}

func (c *channelImpl) GetServer() Server {
	return c.TcpServer
}

func (c *channelImpl) GetMsgChan() chan []byte {
	return c.msgChan
}

func (c *channelImpl) GetMsgBuffChan() chan []byte {
	return c.msgBuffChan
}

//设置链接属性
func (c *channelImpl) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	c.property[key] = value
}

//获取链接属性
func (c *channelImpl) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()
	if value, ok := c.property[key]; ok {
		return value, nil
	} else {
		return nil, errors.New("no property found")
	}
}

//移除链接属性
func (c *channelImpl) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	delete(c.property, key)
}

//发送客户端消息
func SendClientMsg(c Channel, msgId int32, data []byte) error {
	return SendMsg(c, msgId, 0, 0, data)
}

//直接将Message数据发送数据给远程的TCP客户端
func SendMsg(channel Channel, msgId int32, sessionId int64, senderId int64, data []byte) error {
	if channel.IsClose() == true {
		return errors.New("connection closed when send msg")
	}
	//将data封包，并且发送
	var decoder DataPack
	if channel.GetServer().GetServerType() == ClientServer {
		decoder = NewClientDataPack()
		msg, err := decoder.Pack(NewClientMessage(msgId, data))
		if err != nil {
			fmt.Println("Pack error msg id = ", msgId)
			return errors.New("Pack error msg ")
		}
		//写回客户端
		channel.GetMsgChan() <- msg
	} else {
		decoder = NewInnerDataPack()
		msg, err := decoder.Pack(NewInnerMessage(msgId, data, senderId, sessionId))
		if err != nil {
			fmt.Println("Pack error msg id = ", msgId)
			return errors.New("Pack error msg ")
		}
		//写回客户端
		channel.GetMsgChan() <- msg
	}
	return nil
}

//发送proto消息
func SendClientProtoMsg(channel Channel, msgId int32, msg proto.Message) error {
	data, err := proto.Marshal(msg)
	if err != nil {
		fmt.Println("Pack error msg id = ", msgId)
		return errors.New("Pack error msg ")
	}
	return SendMsg(channel, msgId, 0, 0, data)
}

//发送proto消息
func SendProtoMsg(channel Channel, msgId int32, sessionId int64, senderId int64, msg proto.Message) error {
	data, err := proto.Marshal(msg)
	if err != nil {
		fmt.Println("Pack error msg id = ", msgId)
		return errors.New("Pack error msg ")
	}
	return SendMsg(channel, msgId, sessionId, senderId, data)
}

//直接将Message数据发送数据给远程的TCP客户端
func SendBufMsg(channel Channel, msgId int32, sessionId int64, senderId int64, data []byte) error {
	if channel.IsClose() == true {
		return errors.New("connection closed when send msg")
	}
	//将data封包，并且发送
	var decoder DataPack
	if channel.GetServer().GetServerType() == ClientServer {
		decoder = NewClientDataPack()
		msg, err := decoder.Pack(NewClientMessage(msgId, data))
		if err != nil {
			fmt.Println("Pack error msg id = ", msgId)
			return errors.New("Pack error msg ")
		}
		//写回客户端
		channel.GetMsgBuffChan() <- msg
	} else {
		decoder = NewInnerDataPack()
		msg, err := decoder.Pack(NewInnerMessage(msgId, data, senderId, sessionId))
		if err != nil {
			fmt.Println("Pack error msg id = ", msgId)
			return errors.New("Pack error msg ")
		}
		//写回客户端
		channel.GetMsgBuffChan() <- msg
	}
	return nil
}
