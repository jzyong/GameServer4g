package network

import (
	"encoding/binary"
	"fmt"
	"github.com/jzyong/go-mmo-server/src/core/log"
	"io"
	"net"
	"testing"
	"time"
)

// run in terminal:
// go test -v ./znet -run=TestDataPack

//只是负责测试datapack拆包，封包功能
func TestDataPack(t *testing.T) {
	//创建socket TCP Server
	listener, err := net.Listen("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("server listen err:", err)
		return
	}

	//创建服务器gotoutine，负责从客户端goroutine读取粘包的数据，然后进行解析
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("server accept err:", err)
			}

			//处理客户端请求
			go func(conn net.Conn) {

				// 创建拆包解包的对象
				dp := NewClientDataPack()
				// read len
				buffMsgLength := make([]byte, 4)
				if _, err := io.ReadFull(conn, buffMsgLength); err != nil {
					fmt.Println("read msg length error", err)
				}
				var msgLength = uint32(binary.LittleEndian.Uint32(buffMsgLength))
				// 最大长度验证
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
				msg, err := dp.Unpack(msgData, msgLength)
				if err != nil {
					fmt.Println("unpack error ", err)
					return
				}
				fmt.Println("==> Recv Msg: ID=", msg.GetMsgId(), ", len=", msgLength, ", data=", string(msg.GetData()), ",time=", msg.GetTime())

			}(conn)
		}
	}()

	//客户端goroutine，负责模拟粘包的数据，然后进行发送
	go func() {
		conn, err := net.Dial("tcp", "127.0.0.1:7777")
		if err != nil {
			fmt.Println("client dial err:", err)
			return
		}

		//创建一个封包对象 dp
		dp := NewClientDataPack()

		//封装一个msg1包
		msg1 := &ClientMessage{
			Id:   1,
			Time: 123,
			Data: []byte{'h', 'e', 'l', 'l', 'o'},
		}

		sendData1, err := dp.Pack(msg1)
		if err != nil {
			fmt.Println("client pack msg1 err:", err)
			return
		}

		msg2 := &ClientMessage{
			Id:   2,
			Time: 12345,
			Data: []byte{'w', 'o', 'r', 'l', 'd', '!', '!'},
		}
		sendData2, err := dp.Pack(msg2)
		if err != nil {
			fmt.Println("client temp msg2 err:", err)
			return
		}

		//将sendData1，和 sendData2 拼接一起，组成粘包
		sendData1 = append(sendData1, sendData2...)
		//向服务器端写数据
		conn.Write(sendData1)
	}()

	//客户端阻塞
	select {
	case <-time.After(time.Second):
		return
	}
}
