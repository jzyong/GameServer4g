package net

import (
	"encoding/binary"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/jzyong/go-mmo-server/src/core/log"
	network "github.com/jzyong/go-mmo-server/src/core/network/tcp"
	"github.com/jzyong/go-mmo-server/src/message"
	"io"
	"math/rand"
	"net"
	"testing"
	"time"
)

// run in terminal:
// go test -v ./znet -run=TestServer

/*
	模拟客户端
*/
func ClientTest(i int32) {

	fmt.Println("Client Test ... start")
	//3秒之后发起测试请求，给服务端开启服务的机会
	time.Sleep(3 * time.Second)

	conn, err := net.Dial("tcp", "192.168.110.2:6060")
	if err != nil {
		log.Errorf("client start err, exit! %v", err)
		return
	}

	for {
		//发送数据
		dp := network.NewClientDataPack()
		request := &message.ServerListRequest{
			Type: rand.Int31n(100),
		}
		var data, err = proto.Marshal(request)
		msg, _ := dp.Pack(network.NewClientMessage(int32(message.MID_ServerListReq), data))
		_, err = conn.Write(msg)
		if err != nil {
			log.Error("client write err: ", err)
			return
		}

		//接收数据
		dp2 := network.NewClientDataPack()
		// 创建拆包解包的对象
		buffMsgLength := make([]byte, 4)
		// read len
		if _, err := io.ReadFull(conn, buffMsgLength); err != nil {
			log.Error("read msg length error", err)
		}
		var msgLength = uint32(binary.LittleEndian.Uint32(buffMsgLength))
		//最大长度验证
		if msgLength > 10000 {
			log.Warnf("消息太长：%d\n", msgLength)
		}
		msgData := make([]byte, msgLength)
		if _, err := io.ReadFull(conn, msgData); err != nil {
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
		fmt.Println("==> Recv Msg: ID=", msg2.GetMsgId(), ", len=", msgLength, ", data=", string(msg2.GetData()), ",time=", msg2.GetTime())
		response := &message.ServerListResponse{}
		proto.Unmarshal(msg2.GetData(), response)

		log.Infof("收到消息：%v", response)

		time.Sleep(time.Second)
	}
}

func TestConnectClientServer(t *testing.T) {

	//	客户端测试
	go ClientTest(1)

	select {
	case <-time.After(time.Minute * 3):
		return
	}
}
