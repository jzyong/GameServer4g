package net

import (
	"encoding/binary"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/jzyong/go-mmo-server/src/core/log"
	network "github.com/jzyong/go-mmo-server/src/core/network/tcp"
	"github.com/jzyong/go-mmo-server/src/message"
	"io"
	"net"
	"testing"
	"time"
)

// run in terminal:
// go test -v ./znet -run=TestServer

var clientConn net.Conn //客户端连接

/*
	模拟客户端
*/
func ClientTest(i int32) {

	fmt.Println("Client Test ... start")
	//3秒之后发起测试请求，给服务端开启服务的机会
	time.Sleep(3 * time.Second)

	var err error
	clientConn, err = net.Dial("tcp", "192.168.110.2:6060")
	if err != nil {
		log.Errorf("client start err, exit! %v", err)
		return
	}

	go heartRequest()

	go switchReceiveMessage()

}

//分发接收消息
func switchReceiveMessage() {
	for {
		//接收数据
		dp2 := network.NewClientDataPack()
		// 创建拆包解包的对象
		buffMsgLength := make([]byte, 4)
		// read len
		if _, err := io.ReadFull(clientConn, buffMsgLength); err != nil {
			log.Error("read msg length error", err)
		}
		var msgLength = uint32(binary.LittleEndian.Uint32(buffMsgLength))
		//最大长度验证
		if msgLength > 10000 {
			log.Warnf("消息太长：%d\n", msgLength)
		}
		msgData := make([]byte, msgLength)
		if _, err := io.ReadFull(clientConn, msgData); err != nil {
			fmt.Println("read msg data error ", err)
			return
		}
		//fmt.Printf("read headData %+v\n", headData)

		//拆包，得到msgid 和 数据 放在msg中
		msg2, err := dp2.Unpack(msgData, msgLength)
		if err != nil {
			log.Error("unpack error ", err)
			return
		}
		fmt.Println("==> Recv Msg: ID=", msg2.GetMsgId(), ", len=", msgLength)

		switch message.MID(msg2.GetMsgId()) {
		case message.MID_PlayerHeartRes:
			heartResponse(msg2)
		}
		time.Sleep(time.Millisecond * 10)
	}

}

//发送消息 proto 消息
func sendMsg(mid message.MID, message proto.Message) {
	dp := network.NewClientDataPack()
	var data, err = proto.Marshal(message)
	msg, _ := dp.Pack(network.NewClientMessage(int32(mid), data))
	_, err = clientConn.Write(msg)
	if err != nil {
		log.Errorf("消息：%v 发送失败 %v ", mid, err)
		return
	}
}

//发送心跳
func heartRequest() {
	for {
		request := &message.PlayerHeartRequest{}
		sendMsg(message.MID_PlayerHeartReq, request)
		time.Sleep(time.Second * 3)
	}

}

//接收心跳
func heartResponse(tcpMessage network.TcpMessage) {
	response := &message.PlayerHeartResponse{}
	proto.Unmarshal(tcpMessage.GetData(), response)
	log.Infof("收到心跳：%v", response)
}

func TestConnectClientServer(t *testing.T) {

	//	客户端测试
	go ClientTest(1)

	select {
	case <-time.After(time.Minute * 10):
		return
	}
}
