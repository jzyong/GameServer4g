package network

import (
	"bytes"
	"encoding/binary"
)

/*
	封包数据和拆包数据
	直接面向TCP连接中的数据流,为传输数据添加头部信息，用于处理TCP粘包问题。
*/
type DataPack interface {
	//获取包头长度方法
	GetHeadLen() uint32
	//封包方法
	Pack(msg TcpMessage) ([]byte, error)
	//拆包方法
	Unpack([]byte, uint32) (TcpMessage, error)
}

//客户端封包拆包类实例，暂时不需要成员 实现DataPack
type ClientDataPack struct{}

//封包拆包实例初始化方法
func NewClientDataPack() *ClientDataPack {
	return &ClientDataPack{}
}

//获取包头长度方法 ,不包括消息长度
func (dp *ClientDataPack) GetHeadLen() uint32 {
	//Id uint32(4字节) +  timeLength uint64(8字节)
	return 12
}

//封包方法(压缩数据)
func (dp *ClientDataPack) Pack(msg TcpMessage) ([]byte, error) {
	//创建一个存放bytes字节的缓冲
	dataBuff := bytes.NewBuffer([]byte{})

	//写dataLen 不包含自身长度
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetDataLen()+dp.GetHeadLen()); err != nil {
		return nil, err
	}
	//写msgID
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}
	//写时间
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetTime()); err != nil {
		return nil, err
	}
	//写data数据
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}
	return dataBuff.Bytes(), nil
}

//拆包方法(解压数据) 消息长度已经被截取
func (dp *ClientDataPack) Unpack(binaryData []byte, msgLength uint32) (TcpMessage, error) {
	//创建一个从输入二进制数据的ioReader
	dataBuff := bytes.NewReader(binaryData)

	//只解压head的信息，得到dataLen和msgID
	msg := &ClientMessage{}

	//读msgID
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}
	//读时间
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Time); err != nil {
		return nil, err
	}
	//读取数据
	data := make([]byte, msgLength-dp.GetHeadLen())
	if err := binary.Read(dataBuff, binary.LittleEndian, data); err != nil {
		return nil, err
	}
	msg.SetData(data)
	return msg, nil
}

//内部通信封包拆包类实例，暂时不需要成员
type InnerDataPack struct{}

//封包拆包实例初始化方法
func NewInnerDataPack() *InnerDataPack {
	return &InnerDataPack{}
}

//获取包头长度方法 ,不包括消息长度
func (dp *InnerDataPack) GetHeadLen() uint32 {
	//Id uint32(4字节) +  senderId uint64(8字节)+sessionId uint64(8字节)
	return 20
}

//封包方法(压缩数据)
func (dp *InnerDataPack) Pack(msg TcpMessage) ([]byte, error) {
	//创建一个存放bytes字节的缓冲
	dataBuff := bytes.NewBuffer([]byte{})

	//写dataLen 不包含自身长度
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetDataLen()+dp.GetHeadLen()); err != nil {
		return nil, err
	}
	//写msgID
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}
	//发送者id
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetObjectId()); err != nil {
		return nil, err
	}
	//会话id
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetSessionId()); err != nil {
		return nil, err
	}
	//写data数据
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}
	return dataBuff.Bytes(), nil
}

//拆包方法(解压数据) 消息长度已经被截取
func (dp *InnerDataPack) Unpack(binaryData []byte, msgLength uint32) (TcpMessage, error) {
	//创建一个从输入二进制数据的ioReader
	dataBuff := bytes.NewReader(binaryData)
	//只解压head的信息，得到dataLen和msgID
	msg := &InnerMessage{}

	//读msgID
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}
	//发送者id
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.SenderId); err != nil {
		return nil, err
	}
	//会话id
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.SessionId); err != nil {
		return nil, err
	}
	//读取数据
	data := make([]byte, msgLength-dp.GetHeadLen())
	if err := binary.Read(dataBuff, binary.LittleEndian, data); err != nil {
		return nil, err
	}
	msg.SetData(data)
	return msg, nil
}
