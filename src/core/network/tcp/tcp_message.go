package network

import "time"

/*
	将请求的一个消息封装到message中，定义抽象层接口
*/
type TcpMessage interface {
	//获取消息数据段长度
	GetDataLen() uint32
	//获取消息ID
	GetMsgId() int32
	//获取消息内容
	GetData() []byte
	//获取消息时间戳
	GetTime() int64
	//会话id
	GetSessionId() int64
	//对象唯一id
	GetObjectId() int64
	//设置消息ID
	SetMsgId(int32)
	//设置消息内容
	SetData([]byte)
	//设置消息数据段长度
	SetDataLen(uint32)
	//设置时间戳
	SetTime(int64)
	//获取Channel
	GetChannel() Channel
	//设置Channel
	SetChannel(Channel)
}

//内部 消息体 实现TcpMessage
type InnerMessage struct {
	DataLen   uint32  //消息的长度
	Id        int32   //消息的ID
	Data      []byte  //消息的内容
	SessionId int64   //会话id
	SenderId  int64   //发送者id
	Channel   Channel //连接会话
}

//创建一个Message消息包
func NewInnerMessage(id int32, data []byte, senderId int64, sessionId int64) *InnerMessage {
	return &InnerMessage{
		DataLen:   uint32(len(data)),
		Id:        id,
		Data:      data,
		SenderId:  senderId,
		SessionId: sessionId,
	}
}

//获取消息数据段长度
func (msg *InnerMessage) GetDataLen() uint32 {
	return uint32(len(msg.Data))
}

//获取消息ID
func (msg *InnerMessage) GetMsgId() int32 {
	return msg.Id
}

//获取消息内容
func (msg *InnerMessage) GetData() []byte {
	return msg.Data
}

//设置消息数据段长度
func (msg *InnerMessage) SetDataLen(len uint32) {
	msg.DataLen = len
}

//设计消息ID
func (msg *InnerMessage) SetMsgId(msgId int32) {
	msg.Id = msgId
}

//设计消息内容
func (msg *InnerMessage) SetData(data []byte) {
	msg.Data = data
}

//时间
func (msg *InnerMessage) SetTime(time int64) {

}

//时间
func (msg *InnerMessage) GetTime() int64 {
	return 0
}

func (msg *InnerMessage) GetSessionId() int64 {
	return msg.SessionId
}

func (msg *InnerMessage) GetObjectId() int64 {
	return msg.SenderId
}

func (msg *InnerMessage) GetChannel() Channel {
	return msg.Channel
}

func (msg *InnerMessage) SetChannel(channel Channel) {
	msg.Channel = channel
}

//客户端消息体 实现 Message
type ClientMessage struct {
	DataLen uint32  //消息的长度
	Id      int32   //消息的ID
	Data    []byte  //消息的内容
	Time    int64   //消息时间
	Channel Channel //连接会话
}

//创建一个Message消息包
func NewClientMessage(id int32, data []byte) *ClientMessage {
	return &ClientMessage{
		DataLen: uint32(len(data)),
		Id:      id,
		Data:    data,
		Time:    time.Now().Unix() * 1000,
	}
}

//获取消息数据段长度
func (msg *ClientMessage) GetDataLen() uint32 {
	return uint32(len(msg.Data))
}

//获取消息ID
func (msg *ClientMessage) GetMsgId() int32 {
	return msg.Id
}

//获取消息内容
func (msg *ClientMessage) GetData() []byte {
	return msg.Data
}

//设置消息数据段长度
func (msg *ClientMessage) SetDataLen(len uint32) {
	msg.DataLen = len
}

//设计消息ID
func (msg *ClientMessage) SetMsgId(msgId int32) {
	msg.Id = msgId
}

//设计消息内容
func (msg *ClientMessage) SetData(data []byte) {
	msg.Data = data
}

//时间
func (msg *ClientMessage) SetTime(time int64) {
	msg.Time = time
}

//时间
func (msg *ClientMessage) GetTime() int64 {
	return msg.Time
}

func (msg *ClientMessage) GetSessionId() int64 {
	panic("implement me")
}

func (msg *ClientMessage) GetObjectId() int64 {
	panic("implement me")
}

func (msg *ClientMessage) GetChannel() Channel {
	return msg.Channel
}

func (msg *ClientMessage) SetChannel(channel Channel) {
	msg.Channel = channel
}
