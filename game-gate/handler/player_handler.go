package handler

import (
	"github.com/jzyong/GameServer4g/game-message/message"
	"github.com/jzyong/golib/log"
	network "github.com/jzyong/golib/network/tcp"
	"time"
)

//func init() {
//	manager.GetClientManager().RegisterHandler(int32(message.MID_PlayerHeartReq), HandlePlayerHeartReq)
//}

//玩家心跳请求
func HandlePlayerHeartReq(msg network.TcpMessage) bool {
	response := &message.PlayerHeartResponse{
		Timestamp: time.Now().Unix(),
	}
	log.Info("%v 返回心跳%d", msg.GetChannel().RemoteAddr(), response.GetTimestamp())
	network.SendClientProtoMsg(msg.GetChannel(), int32(message.MID_PlayerHeartRes), response)
	return true
}
