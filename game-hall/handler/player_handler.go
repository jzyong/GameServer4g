package handler

import (
	"context"
	"github.com/golang/protobuf/proto"
	"github.com/jzyong/GameServer4g/game-hall/manager"
	"github.com/jzyong/GameServer4g/game-message/message"
	"github.com/jzyong/golib/log"
	network "github.com/jzyong/golib/network/tcp"
	"github.com/jzyong/golib/util"
	"time"
)

//func init() {
//	manager.GetClientManager().MessageDistribute.RegisterHandler(int32(message.MID_UserLoginReq), network.NewTcpHandler(HandUserLogin))
//}

//处理玩家登录
func HandUserLogin(msg network.TcpMessage) bool {
	request := &message.UserLoginRequest{}
	proto.Unmarshal(msg.GetData(), request)
	log.Info("请求账号：%v 密码：%v", request.GetAccount(), request.GetPassword())

	//TODO 添加MongoDB存储
	playerId, _ := util.UUID.GetId()
	response := &message.UserLoginResponse{
		PlayerId: playerId,
	}

	//登录世界服，测试
	c, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	//TODO nil 处理
	worldResponse, _ := manager.GetClientManager().PlayerWorldClient.Login(c, request)
	log.Info("world return %v", worldResponse)

	manager.SendMsg(msg.GetChannel(), int32(message.MID_UserLoginRes), playerId, response)
	return true
}
