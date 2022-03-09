package rpc

import (
	"context"
	"github.com/jzyong/GameServer4g/game-message/message"
	"github.com/jzyong/golib/log"
)

type PlayerServiceImpl struct {
}

//玩家进入世界服
func (service *PlayerServiceImpl) Login(ctx context.Context, request *message.UserLoginRequest) (*message.UserLoginResponse, error) {
	log.Info("%s 登录", request.GetAccount())
	response := &message.UserLoginResponse{
		PlayerId: 1,
	}
	return response, nil
}
