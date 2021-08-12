package handler

import (
	"context"
	"github.com/golang/protobuf/proto"
	"github.com/jzyong/go-mmo-server/src/core/log"
	network "github.com/jzyong/go-mmo-server/src/core/network/tcp"
	"github.com/jzyong/go-mmo-server/src/core/util"
	"github.com/jzyong/go-mmo-server/src/hall/manager"
	"github.com/jzyong/go-mmo-server/src/message"
	"time"
)

func init() {
	manager.GetClientManager().MessageDistribute.RegisterHandler(int32(message.MID_UserLoginReq), network.NewTcpHandler(HandUserLogin))
}

//处理玩家登录
func HandUserLogin(msg network.TcpMessage) bool {
	request := &message.UserLoginRequest{}
	proto.Unmarshal(msg.GetData(), request)
	log.Infof("请求账号：%v 密码：%v", request.GetAccount(), request.GetPassword())

	//TODO 添加MongoDB存储
	playerId, _ := util.UUID.GetId()
	response := &message.UserLoginResponse{
		PlayerId: playerId,
	}

	//登录世界服，测试
	c, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	worldResponse, _ := manager.GetClientManager().PlayerWorldClient.Login(c, request)
	log.Infof("world return %v", worldResponse)

	manager.SendMsg(msg.GetChannel(), int32(message.MID_UserLoginRes), playerId, response)
	return true
}
