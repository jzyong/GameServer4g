package rpc

import (
	"context"
	"github.com/jzyong/go-mmo-server/src/core/log"
	"github.com/jzyong/go-mmo-server/src/message"
)

type PlayerServiceImpl struct {
}

//玩家进入世界服
func (service *PlayerServiceImpl) Login(ctx context.Context, request *message.UserLoginRequest) (*message.UserLoginResponse, error) {
	log.Infof("%s 登录", request.GetAccount())
	response := &message.UserLoginResponse{
		PlayerId: 1,
	}
	return response, nil
}
